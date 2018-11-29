# Cherwell Service Manager Client

Go client for Cherwell Service Management (CSM) 

## Documentation
For details on functionality, see [Cherwell API Golang library](https://continuum.atlassian.net/wiki/spaces/C2E/pages/781945314/Cherwell+API+Golang+library)

### Setup
How to create a client:
```go
import "github.com/ContinuumLLC/platform-common-lib/src/cherwell"
        
var conf = cherwell.Config{
                        Host    : "CSM_HOST",
                        ClientId: "CSM_CLIENT_ID",
                        UserName: "CSM_USERNAME",
                        PassWord: "CSM_PASSWORD"}

cl, err := cherwell.NewClient(conf, &http.Client{})
 	if err != nil {
 		// handle error
 	}
```
#### Reading Business Objects from CSM

To get information about Business Objects from CSM use following methods:

- `func (c *Client) GetByPublicID(busobID, publicID string) (*BusinessObject, error)` - 
method that returns a Business Object record that includes a list of fields and their record IDs, names, and set values.
Specify the Business Object ID and Business Object public ID.
- `func (c *Client) GetByRecordID(busobID, busobRecID string) (*BusinessObject, error)` - 
method that returns a Business Object record that includes a list of fields and their record IDs, 
names, and set values. Specify the Business Object ID and Business Object record ID.
- `func (c *Client) GetBatch(items []BusinessObjectInfo) ([]BusinessObject, error)` - 
method that returns a batch of Business Object records that includes a list of field record IDs, display names, 
and values for each record. Specify an array of Business Object IDs, record IDs, or public IDs ([]BusinessObjectInfo). 
Use a flag to stop on error or continue on error.

#### Creating Business Objects in CSM

To create Business Objects in CSM use following methods:

- `func (c *Client) Save(bo BusinessObject) (string, string, error)` - to create Business Object, 
leave record ID and public ID empty. 
Upon creating a cache key is returned to use for subsequent requests. If the object is not found in the cache 
with said cache key, specify record ID or public ID to save and return a new cache key. Set persist = true, to actually 
save the Business Object to disk, persist = false will just cache it.
- `func (c *Client) SaveBatch(bos []BusinessObject) ([]string, []string, error)` - to create an array of Business 
Objects in a batch, leave record ID and public ID empty.

#### Updating Business Objects in CSM

To update Business Objects in CSM use the same methods like for creating:

- `func (c *Client) Save(bo BusinessObject) (string, string, error)` - to update Business Object specify
Business Object ID and record ID or public ID.
- `func (c *Client) SaveBatch(bos []BusinessObject) ([]string, []string, error)` - to update Business Objects 
specify Business Object IDs and record IDs or public IDs specify record ID or public ID.

#### Deleting Business Objects from CSM

To delete Business Objects from CSM use following methods:

- `func (c *Client) DeleteByPublicID(busobID, publicID string) (string, string, error)` - 
method that deletes a Business Object record by public ID. Specify the Business Object ID and Business Object public ID.
- `func (c *Client) DeleteByRecordID(busobID, busobRecID string) (string, string, error)` - 
method that deletes a Business Object record by Record ID. Specify the Business Object ID and Business Object record ID.
- `func (c *Client) DeleteBatch(itemsToDelete []BusinessObjectInfo) ([]string, []string, error)` - 
method that deletes a batch of Business Object records. Specify an array of Business Object IDs, record IDs
or public IDs. Use a flag to stop on error or continue on error.

### Examples

This library supports CRUD operations on Business Objects.

Create Business Object:
```go
    /* bo - variable which stores information about BusinessObject:
        type BusinessObject struct {
        	BusObID       string              `json:"busObId"`
        	BusObPublicID string              `json:"busObPublicId"`
        	BusObRecID    string              `json:"busObRecId"`
        	Fields        []FieldTemplateItem `json:"fields"`
        }
    */
    busObPublicID, busObRecID, err := cl.Save(bo)
    if err != nil { 
            // handle error
    }	
```
Create Business Objects in batch:
```go
    /* bos - variable which stores information about BusinessObjects:
        var bos []BusinessObject
    */
    publicIds, recIds, err := cl.SaveBatch(bos)
    if err != nil {
        // handle error
    }
```
Read Business Object by Public ID:
```go
    bo, err := cl.GetByPublicID(busobID, publicID)
    if err != nil {
		// handle error
    }
```
Read Business Object by Record ID:
```go
    bo, err := cl.GetByRecordID(busobID, recID)
    if err != nil {
		// handle error
    }
```

Read Business Objects in batch:
```go
    /* type BusinessObjectInfo struct {
            BusObID       string `json:"busObId"`
            BusObPublicID string `json:"busObPublicId"`
            BusObRecID    string `json:"busObRecId"`
       }
    */
    var boInfo []BusinessObjectInfo
    // set up particular boInfo
    bos, err := cl.GetBatch(boInfo)
    if err != nil {
		// handle error
    }
```

Delete Business Object by Public ID:
```go
    bo, err := cl.DeleteByPublicID(busobID, publicID)
    if err != nil {
		// handle error
    }
```

Delete Business Object by Record ID:
```go
    bo, err := cl.DeleteByRecordID(busobID, recID)
    if err != nil {
		// handle error
    }
```

Delete Business Objects in batch:
```go
    var itemsToDelete []BusinessObjectInfo
    pubIDs, recIDs, err := cl.DeleteBatch(itemsToDelete)
    if err != nil {
		// handle error
    }
```

Update Business Object (the same method as for saving of Business Object):
```go
    /* bo - variable which stores information about BusinessObject with fields to update:
        type BusinessObject struct {
        	BusObID       string              `json:"busObId"`
        	BusObPublicID string              `json:"busObPublicId"`
        	BusObRecID    string              `json:"busObRecId"`
        	Fields        []FieldTemplateItem `json:"fields"`
        }
    */
    busObPublicID, busObRecID, err := cl.Save(bo)
    if err != nil { 
            // handle error
    }	
```

Update Business Objects in batch (the same method as for saving of Business Objects):
```go
    /* bos - variable which stores information about BusinessObjects with fields to update:
        var bos []BusinessObject
    */
    publicIds, recIds, err := cl.SaveBatch(bos)
    if err != nil {
        // handle error
    }
```