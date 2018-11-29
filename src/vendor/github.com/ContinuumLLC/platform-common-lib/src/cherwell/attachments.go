package cherwell

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

var attachmentFileNameRegEx = regexp.MustCompile(`(?m)filename=["]?([^"*?]+)["]?`)

// RecordAttachmentType is record attachment type
type RecordAttachmentType string

// AttachmentType is file attachment type.
type AttachmentType string

/////////////////////////////
// Record attachment types
/////////////////////////////
const (
	// FileRecord is linked files record attachment type
	FileRecord RecordAttachmentType = "File"

	// FileManagerFileRecord is imported files record attachment type
	FileManagerFileRecord RecordAttachmentType = "FileManagerFile"

	// BusObRecord is attached business objects record attachment type
	BusObRecord RecordAttachmentType = "BusOb"

	// History is attached history record attachment type
	HistoryRecord RecordAttachmentType = "History"
)

/////////////////////////////
// Attachment types
/////////////////////////////
const (
	// ImportedAttachment type means that attachment was imported into database
	ImportedAttachment AttachmentType = "Imported"

	// LinkedAttachment type means that attachment is linked to an external file
	LinkedAttachment AttachmentType = "Linked"

	// URLAttachment type means that attachment is URL
	URLAttachment AttachmentType = "URL"
)

// AttachmentsQuery is query for selecting attachments
type AttachmentsQuery struct {
	AttachmentID    string   `json:"attachmentId,omitempty"`
	AttachmentTypes []string `json:"attachmentTypes,omitempty"`
	ObjectID        string   `json:"busObId"`
	PublicID        string   `json:"busObPublicId,omitempty"`
	RecordID        string   `json:"busObRecId,omitempty"`
	IncludeLinks    bool     `json:"includeLinks"`
	Types           []string `json:"types,omitempty"`
}

// AttachmentResponse is attachments query response
type AttachmentResponse struct {
	ErrorData
	Attachments []AttachmentSummary `json:"attachments"`
}

// AttachedFile is downloaded attachment file
type AttachedFile struct {
	FileName    string
	ContentType string
	SizeBytes   string
	Data        io.ReadCloser
}

// AttachmentSummary contains information about attachment
type AttachmentSummary struct {
	BusinessObjectInfo
	FileID         string `json:"attachmentFileId"`
	FileName       string `json:"attachmentFileName"`
	FileType       string `json:"attachmentFileType"`
	AttachmentID   string `json:"attachmentId"`
	AttachmentType int    `json:"attachmentType"`
	Comment        string `json:"comment"`
	CreatedAt      string `json:"created"`
	DisplayText    string `json:"displayText"`
	Owner          string `json:"owner"`
	Scope          int    `json:"scope"`
	ScopeOwner     string `json:"scopeOwner"`
	Type           int    `json:"type"`
	Links          []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"links"`
}

// Attachment represents Cherwell BO attachment
type Attachment struct {
	Owner        BusinessObjectInfo
	FileName     string
	Description  string
	Data         io.Reader
	Offset       int
	AttachmentID string
}

// getUploadParams returns attachment contents as buffer, relative upload path with file metadata and error (if any)
func (a *Attachment) getUploadParams() (buffer io.Reader, path string, err error) {
	buf := new(bytes.Buffer)
	bytesRead, err := buf.ReadFrom(a.Data)
	if err != nil {
		err = fmt.Errorf("failed to read attachment data from the buffer, %v", err)
		return
	}

	displayText := a.Description
	if displayText == "" {
		displayText = a.FileName
	}

	path = fmt.Sprintf(
		attachmentUploadPathTemplate+"?displaytext=%s",
		a.FileName,
		a.Owner.ID,
		a.Owner.RecordID,
		a.Offset,
		strconv.FormatInt(bytesRead, 10),
		displayText,
	)

	if a.AttachmentID != "" {
		path += "&attachmentid=" + a.AttachmentID
	}

	buffer = buf
	return
}

// NewAttachment creates a new attachment object
func NewAttachment(fileName string, data io.Reader, owner *BusinessObjectInfo) *Attachment {
	return &Attachment{
		Owner:    *owner,
		FileName: fileName,
		Data:     data,
	}
}

// UploadAttachment uploads a file as attachment to business object
func (c *Client) UploadAttachment(a *Attachment) (fileID string, err error) {
	buff, absPath, err := a.getUploadParams()
	if err != nil {
		return
	}

	reqPath := attachmentUploadPath + absPath
	req, err := c.createRequest(http.MethodPost, reqPath, buff)
	if err != nil {
		err = fmt.Errorf("failed to create upload request, %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to upload attachment: %v", err.Error())
		return
	}

	return c.readUploadResponse(resp)
}

// GetObjectAttachments gets attachments for specific business object
func (c *Client) GetObjectAttachments(boid, recordID string, recordType RecordAttachmentType, attachmentType AttachmentType) (resp *AttachmentResponse, err error) {
	resp = new(AttachmentResponse)
	reqPath := fmt.Sprintf(attachmentGetPath, boid, recordID, recordType, attachmentType)
	err = c.performRequest(http.MethodGet, reqPath, nil, resp)
	if err != nil {
		return
	}

	if resp.HasError {
		err = resp.GetErrorObject()
		return
	}

	return
}

// DeleteAttachment removes an attachment
func (c *Client) DeleteAttachment(attachmentID string, owner *BusinessObjectInfo) error {
	reqPath := fmt.Sprintf(
		attachmentDeletePath,
		attachmentID,
		owner.ID,
		owner.RecordID,
	)
	resp := &ErrorData{}
	err := c.performRequest(http.MethodDelete, reqPath, nil, resp)
	if err != nil {
		return err
	}

	if resp.HasError {
		return resp.GetErrorObject()
	}

	return nil
}

// extractAttachmentFileName extracts file name from Content-Disposition header.
//
// Value sample: 'inline; filename="Some file.docx"'
func (c *Client) extractAttachmentFileName(contentDispositionValue string) (string, error) {
	matches := attachmentFileNameRegEx.FindAllStringSubmatch(contentDispositionValue, -1)

	// Filename should be located at first match at first group
	if len(matches) == 0 || len(matches[0]) < 2 {
		return "", fmt.Errorf("file name not found in '%s'", contentDispositionValue)
	}

	return matches[0][1], nil
}

// AttachmentByID gets attachment contents by ID
func (c *Client) AttachmentByID(attachmentID string, owner *BusinessObjectInfo) (*AttachedFile, error) {
	reqPath := fmt.Sprintf(
		attachmentDownloadPath,
		attachmentID,
		owner.ID,
		owner.RecordID,
	)

	resp, err := c.getResponse(http.MethodGet, reqPath, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close() // nolint: errcheck
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// Return HTTP status as error,
			// if response error body is unreadable
			return nil, fmt.Errorf(resp.Status)
		}

		return nil, errorFromResponse(string(body))
	}

	fileName, err := c.extractAttachmentFileName(
		resp.Header.Get("Content-Disposition"),
	)

	if err != nil {
		// Set default file name, if file name is not available
		fileName = attachmentID + ".txt"
	}

	return &AttachedFile{
		FileName:    fileName,
		ContentType: resp.Header.Get("Content-Type"),
		SizeBytes:   resp.Header.Get("Content-Length"),
		Data:        resp.Body,
	}, nil
}

func (c *Client) readUploadResponse(resp *http.Response) (fileID string, err error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response body from Cherwell, %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errorFromResponse(string(body))
		return
	}

	fileID = string(body)
	return
}
