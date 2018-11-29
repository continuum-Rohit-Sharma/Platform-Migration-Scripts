package cherwell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestAttachmentFileNameParse(t *testing.T) {
	c := Client{}
	cases := map[string]string{
		"file-without-quotes.png":  `inline; filename=file-without-quotes.png`,
		"file-with-quotes.png":     `inline; filename="file-with-quotes.png"`,
		"file with spaces (1).png": `inline; filename="file with spaces (1).png"`,
	}

	for expected, input := range cases {
		got, err := c.extractAttachmentFileName(input)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, got)
	}
}

func TestAttachmentDownload(t *testing.T) {
	cases := map[string]struct {
		id                 string
		owner              *BusinessObjectInfo
		fileName           string
		contentDisposition string
		expError           string
		error              *ErrorData
	}{
		"should download valid attachment": {
			id:       "1",
			fileName: "./testdata/favicon.ico",
			owner: &BusinessObjectInfo{
				ID:       "2",
				RecordID: "200",
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(tt *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					tt.Fatal(r)
				}
			}()

			var expected *AttachedFile
			shouldError := test.expError != ""
			cdOverrided := test.contentDisposition != ""

			var expPayload []byte
			var buff bytes.Buffer

			if test.fileName != "" {
				f, err := os.Open(test.fileName)
				if err != nil {
					tt.Fatal(err)
				}

				info, err := f.Stat()
				if err != nil {
					tt.Fatal(err)
				}

				expected = &AttachedFile{
					FileName:    filepath.Base(info.Name()),
					ContentType: "application/octet-stream",
					SizeBytes:   strconv.FormatInt(info.Size(), 10),
					Data:        f,
				}

				// TeeReader is needed to read expected Data twice:
				// from http handler and from assert function
				reader := io.TeeReader(expected.Data, &buff)
				expPayload, _ = ioutil.ReadAll(reader)
				defer f.Close()
			}

			reqPath := fmt.Sprintf(
				attachmentDownloadPath,
				test.id,
				test.owner.ID,
				test.owner.RecordID,
			)

			server, mux := newTestServer()
			defer func() {
				server.CloseClientConnections()
				server.Close()
			}()

			mux.Handle(
				reqPath,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if test.error != nil {
						w.WriteHeader(500)
						r, _ := json.Marshal(test.error)
						w.Write(r)
						return
					}

					var cdHeader string
					if cdOverrided {
						cdHeader = test.contentDisposition
					} else {
						cdHeader = fmt.Sprintf(`attachment; filename="%s"`, expected.FileName)
					}

					w.WriteHeader(200)
					headers := w.Header()
					headers.Set("Content-Type", expected.ContentType)
					headers.Set("Content-Disposition", cdHeader)
					headers.Set("Content-Length", expected.SizeBytes)
					written, err := io.Copy(w, &buff)
					if err != nil {
						tt.Fatalf("failed to write response payload: %v", err)
					}
					t.Logf("%d bytes written", written)
					return
				}),
			)

			client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
			if err != nil {
				tt.Fatal(err)
			}

			result, err := client.AttachmentByID(test.id, test.owner)
			if err != nil {
				if !shouldError {
					tt.Fatal(err)
				}

				if !strings.Contains(err.Error(), test.expError) {
					t.Fatalf("error message '%s' should contain '%s'", err, test.expError)
				}
			}

			gotPayload, _ := ioutil.ReadAll(result.Data)
			assert.ElementsMatchf(tt, expPayload, gotPayload, "Expected and received data mismatch")
		})
	}
}

func TestAttachmentUpload(t *testing.T) {
	cases := map[string]struct {
		fileName   string
		bo         BusinessObjectInfo
		expError   string
		expCode    int
		expectedID string
	}{
		"should return asset ID on upload success": {
			fileName:   "image.jpg",
			expCode:    200,
			expectedID: "94412c430b734b58c7f3004979be0f6e089160455d",
			bo: BusinessObjectInfo{
				ID:       "6dd53665c0c24cab86870a21cf6434ae",
				RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
			},
		},
		"should return proper error from non-json response": {
			fileName: "image.jpg",
			expCode:  500,
			expError: "Here Is Non JSON Error",
			bo: BusinessObjectInfo{
				ID:       "6dd53665c0c24cab86870a21cf6434ae",
				RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
			},
		},
	}

	for tName, tCase := range cases {
		t.Run(tName, func(tt *testing.T) {
			fPath := "./testdata/" + tCase.fileName
			fileReader, err := os.Open(fPath)
			if err != nil {
				tt.Fatal(err)
			}

			fi, e := fileReader.Stat()
			if e != nil {
				tt.Fatal(e)
			}

			a := NewAttachment(tCase.fileName, fileReader, &tCase.bo)
			server, mux := newTestServer()
			client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
			if err != nil {
				tt.Fatal(err)
			}

			defer func() {
				server.CloseClientConnections()
				server.Close()
			}()

			expPath := fmt.Sprintf(
				attachmentUploadPathTemplate,
				tCase.fileName,
				tCase.bo.ID,
				tCase.bo.RecordID,
				0,
				strconv.FormatInt(fi.Size(), 10),
			)

			output := tCase.expectedID
			shouldError := tCase.expError != ""
			if shouldError {
				output = tCase.expError
			}

			mux.Handle(
				attachmentUploadPath+expPath,
				newMockHandler(http.MethodPost, attachmentUploadPath+expPath, output, tCase.expCode),
			)

			result, err := client.UploadAttachment(a)
			if err != nil {
				if shouldError && strings.Contains(err.Error(), tCase.expError) {
					tt.Logf("%s - OK", tName)
					return
				} else {
					tt.Fatalf("Unexpected error message:\n  Got: %s\n  Exp: %s\n", err.Error(), tCase.expError)
				}
			}

			if result != tCase.expectedID {
				tt.Fatalf("ID mismatch (%s != %s)", tCase.expectedID, result)
			}
		})
	}
}

func TestAttachmentDelete(t *testing.T) {
	cases := map[string]struct {
		atID        string
		owner       BusinessObjectInfo
		code        int
		shouldError bool
		expError    ErrorData
		httpError   string
	}{
		"should return ok response on delete ok": {
			atID: "94412c430b734b58c7f3004979be0f6e089160455d",
			code: 200,
			owner: BusinessObjectInfo{
				ID:       "6dd53665c0c24cab86870a21cf6434ae",
				RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
			},
		},
		"should return correct error on non-json response": {
			atID:      "123",
			code:      500,
			httpError: "something went bad",
			owner: BusinessObjectInfo{
				ID:       "6dd53665c0c24cab86870a21cf6434ae",
				RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
			},
		},
		"should return correct error on error response": {
			atID:        "234",
			code:        500,
			shouldError: true,
			owner: BusinessObjectInfo{
				ID:       "6dd53665c0c24cab86870a21cf6434ae",
				RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
			},
			expError: ErrorData{
				HasError:     true,
				ErrorCode:    BusObNotValidError,
				ErrorMessage: "Business object is not valid",
			},
		},
	}

	for name, k := range cases {
		t.Run(name, func(tt *testing.T) {
			server, mux := newTestServer()
			client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
			if err != nil {
				tt.Fatal(err)
			}

			defer func() {
				server.CloseClientConnections()
				server.Close()
			}()

			if err != nil {
				t.Fatalf("can't prepare output, %v", err)
			}

			output, err := json.Marshal(k.expError)
			if err != nil {
				tt.Fatal(err)
			}

			reqPath := fmt.Sprintf(
				attachmentDeletePath,
				k.atID,
				k.owner.ID,
				k.owner.RecordID,
			)
			mux.Handle(
				reqPath,
				newMockHandler(http.MethodDelete, "-", string(output), k.code),
			)

			err = client.DeleteAttachment(k.atID, &k.owner)
			if err != nil {
				var expectedError string
				if k.httpError != "" {
					expectedError = k.httpError
				} else {
					expectedError = k.expError.ErrorMessage
				}

				if k.shouldError && strings.Contains(err.Error(), expectedError) {
					return
				}

				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAttachmentGet(t *testing.T) {
	cases := map[string]struct {
		atID      string
		owner     BusinessObjectInfo
		code      int
		expError  ErrorData
		httpError string
		resp      *AttachmentResponse
	}{
		"should return ok response on delete ok": {
			atID: "94412c430b734b58c7f3004979be0f6e089160455d",
			code: 200,
			owner: BusinessObjectInfo{
				ID:       "6dd53665c0c24cab86870a21cf6434ae",
				RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
			},
			resp: &AttachmentResponse{
				Attachments: []AttachmentSummary{
					{
						BusinessObjectInfo: BusinessObjectInfo{
							RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
						},
					},
				},
			},
		},
		"should return correct error on non-json response": {
			atID:      "123",
			code:      500,
			httpError: "something went bad",
			owner: BusinessObjectInfo{
				ID:       "6dd53665c0c24cab86870a21cf6434ae",
				RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
			},
		},
		"should return correct error on error response": {
			atID: "234",
			code: 500,
			owner: BusinessObjectInfo{
				ID:       "6dd53665c0c24cab86870a21cf6434ae",
				RecordID: "94412acbaab57e55fe60df4d5688591d9c63544377",
			},
			expError: ErrorData{
				HasError:     true,
				ErrorCode:    BusObNotValidError,
				ErrorMessage: "Business object is not valid",
			},
		},
	}

	for name, k := range cases {
		t.Run(name, func(tt *testing.T) {
			server, mux := newTestServer()
			client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
			if err != nil {
				tt.Fatal(err)
			}

			defer func() {
				server.CloseClientConnections()
				server.Close()
			}()

			var output []byte
			shouldError := false
			if k.resp != nil {
				output, err = json.Marshal(k.resp)
			} else if k.httpError != "" {
				shouldError = true
				output = []byte(k.httpError)
			} else {
				shouldError = true
				output, err = json.Marshal(k.expError)
			}

			if err != nil {
				t.Fatalf("can't prepare output, %v", err)
			}

			reqPath := fmt.Sprintf(attachmentGetPath, k.owner.ID, k.owner.RecordID, FileRecord, ImportedAttachment)
			mux.Handle(
				reqPath,
				newMockHandler(http.MethodGet, reqPath, string(output), k.code),
			)

			resp, err := client.GetObjectAttachments(
				k.owner.ID,
				k.owner.RecordID,
				FileRecord,
				ImportedAttachment,
			)

			if err != nil {
				var expectedError string
				if k.httpError != "" {
					expectedError = k.httpError
				} else {
					expectedError = k.expError.ErrorMessage
				}

				if shouldError && strings.Contains(err.Error(), expectedError) {
					tt.Logf("%s - OK", name)
					return
				} else {
					tt.Fatalf("Unexpected error message:\n  Got: %s\n  Exp: %s\n", err.Error(), expectedError)
				}
			}

			if k.resp == nil {
				return
			}

			assert.ElementsMatchf(tt, resp.Attachments, k.resp.Attachments, "Got:\n  %#v\nWant:\n  %#v", resp.Attachments, k.resp.Attachments)
		})
	}
}
