package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"os"
	"reflect"
)

//KEYTOSEARCH is used to search of key for replacement in schedule
const KEYTOSEARCH string = "name"

func newMergeService() MergeService {
	return &mergeServiceImpl{}
}

type mergeServiceImpl struct{}

func (m *mergeServiceImpl) Merge(cfg *Config) (err error) {
	mapSrcObj := make(map[string]interface{})
	mapSrc := &mapSrcObj
	mapDstObj := make(map[string]interface{})
	mapDst := &mapDstObj

	switch cfg.Action {
	case "remove":
		cfg.mode = &modeLess{}
	case "equal":
		cfg.mode = &modeSame{}
	default:
		cfg.mode = &modeAdd{}
	}

	err = m.parseFile(cfg.Source, mapSrc)
	if err != nil {
		return
	}
	//fmt.Printf("%+v\n", maps)

	err = m.parseFile(cfg.Delta, mapDst)
	if err != nil {
		return
	}
	err = m.mergeMaps(mapSrc, mapDst, cfg)
	if err != nil {
		return
	}

	err = m.dumpFile(cfg.Destination, mapSrc)
	if err != nil {
		return
	}

	return
}

func (m *mergeServiceImpl) mergeMaps(mapSrc, mapDst *map[string]interface{}, cfg *Config) (err error) {
	mapDstObj := *mapDst
	mapSrcObj := *mapSrc
	for k, v := range mapDstObj {
		srcValue, srcHasKey := mapSrcObj[k]
		if !srcHasKey {
			cfg.mode.SourceNotHasKey(mapSrc, mapDst, k, v)
		} else {
			rt := reflect.TypeOf(srcValue)
			switch rt.Kind() {
			case reflect.Map:
				intrt := reflect.TypeOf(v)
				switch intrt.Kind() {
				case reflect.Map:
					prmSrc, _ := srcValue.(map[string]interface{})
					prmDst, _ := v.(map[string]interface{})
					err = m.mergeMaps(&prmSrc, &prmDst, cfg)
					if err != nil {
						return errors.New(err.Error() + ":" + k)
					}
				// case reflect.Array:
				// 	cfg.mode.IfArray(mapSrc, mapDst, k)
				default:
					return errors.New("TypeMismatch: " + k)
				}
			case reflect.Array:
			case reflect.Slice:
				// get array-item-prop d dst
				var (
					itemPropStr string
					itemValue   interface{}
				)

				switch reflect.TypeOf(v).Kind() {
				case reflect.Map:
					{
						var itemValueIsDefined bool
						dstParentMap := v.(map[string]interface{})
						itemProp, itemPropIsDefined := dstParentMap["prop"]
						itemValue, itemValueIsDefined = dstParentMap["value"]
						if itemPropIsDefined && itemValueIsDefined &&
							reflect.TypeOf(itemProp).Kind() == reflect.String &&
							(reflect.TypeOf(itemValue).Kind() == reflect.Array ||
								reflect.TypeOf(itemValue).Kind() == reflect.Slice) {
							// construct src map of array itemsd
							itemPropStr = itemProp.(string)

						}

					}
				case reflect.Array:
				case reflect.Slice:
					{
						itemPropStr = KEYTOSEARCH
						itemValue = v
					}
				default:
					{
						return errors.New("Not supported format for delta file")
					}
				}
				srcItemsMap := make(map[string]interface{})
				srcValueArr := srcValue.([]interface{})
				for _, srcItem := range srcValueArr {
					srcItemMap := srcItem.(map[string]interface{})
					srcItemsMap[srcItemMap[itemPropStr].(string)] = srcItem
				}
				// construct dst map of array items
				dstItemsMap := make(map[string]interface{})
				dstValueArr := itemValue.([]interface{})
				for _, dstItem := range dstValueArr {
					dstItemMap := dstItem.(map[string]interface{})
					dstItemsMap[dstItemMap[itemPropStr].(string)] = dstItem
				}

				m.mergeMaps(&srcItemsMap, &dstItemsMap, cfg)

				mapSrcObj[k] = m.arrayOf(&srcItemsMap)

			//TODO - what to do if my special delta array handling is not defined?

			default:
				{
					cfg.mode.SourceHasKey(mapSrc, mapDst, k, v)
				}
			}
		}
	}
	return
}

func (m *mergeServiceImpl) arrayOf(mp *map[string]interface{}) interface{} {
	arr := make([]interface{}, 0, len(*mp))
	for _, value := range *mp {
		if reflect.ValueOf(value).Len() > 0 {
			arr = append(arr, value)
		}
	}
	return arr
}

func (m *mergeServiceImpl) dumpFile(filePath string, maps *map[string]interface{}) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	ser := json.NewEncoder(file)
	ser.SetIndent("", "\t")
	err = ser.Encode(maps)
	if err != nil {
		return
	}

	return
}

func (m *mergeServiceImpl) parseFile(filePath string, maps *map[string]interface{}) (err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	deser := json.NewDecoder(file)
	err = deser.Decode(maps)
	if err != nil {
		return
	}

	return
}

type modeAdd struct{}

func (m *modeAdd) SourceHasKey(src, dest *map[string]interface{}, key string, value interface{}) {
	mapSrcObj := *src
	mapSrcObj[key] = value
}

func (m *modeAdd) SourceNotHasKey(src, dest *map[string]interface{}, key string, value interface{}) {
	mapSrcObj := *src
	mapSrcObj[key] = value
}

func (m *modeAdd) IfArray(src, dest *map[string]interface{}, key string) {
	// mapSrcObj := *src
	// mapDestObj := *dest
	// ms := mergeServiceImpl{}
	// ms.mergeMaps(&mapSrcObj, &mapDestObj, nil)
	// for k, v := range mapDestObj {
	// 	mapDestObj[k]
	// }

	// get array-config { prop: __, value: []} from dst
	// get prop
	// construct srcIndexMap
	// // new map map[interface{}]interface{}
	// // loop srcArray, populate srcIndexMap
	// construct dstIndexMap
	// // as above with dstArray
	// construct finalIndexMap - clone of srcIndexMap
	// call mergeMaps(srcIndexMap)
}

type modeLess struct{}

func (m *modeLess) SourceHasKey(src, dest *map[string]interface{}, key string, value interface{}) {
	delete(*src, key)
}

func (m *modeLess) SourceNotHasKey(src, dest *map[string]interface{}, key string, value interface{}) {
	// do nothing
}

func (m *modeLess) IfArray(src, dest *map[string]interface{}, key string) {
}

type modeSame struct{}

func (m *modeSame) SourceHasKey(src, dest *map[string]interface{}, key string, value interface{}) {
	mapSrcObj := *src
	mapSrcObj[key] = value
}

func (m *modeSame) SourceNotHasKey(src, dest *map[string]interface{}, key string, value interface{}) {
	// do nothing
}

func (m *modeSame) IfArray(src, dest *map[string]interface{}, key string) {
	// mapSrcObj := *src
	// mapDestObj := *dest
	// srcValue, _ := mapSrcObj[key]
	// destValue, _ := mapDestObj[key]
	// srcObj, destObj := createArrayObj(srcValue, destValue)
}

func createArrayObj(src, dest interface{}) (srcObj, destObj *Schedules) {
	srcSchedules := Schedules{}
	destSchedules := Schedules{}
	var srcbuf bytes.Buffer
	var destbuf bytes.Buffer
	srcenc := gob.NewEncoder(&srcbuf)
	err := srcenc.Encode(src)
	if err != nil {
		return nil, nil
	}
	err = json.Unmarshal(srcbuf.Bytes(), &srcSchedules)

	destenc := gob.NewEncoder(&destbuf)
	err = destenc.Encode(dest)
	if err != nil {
		return nil, nil
	}
	err = json.Unmarshal(srcbuf.Bytes(), &destSchedules)
	if err != nil {
		return nil, nil
	}

	return &srcSchedules, &destSchedules
}
