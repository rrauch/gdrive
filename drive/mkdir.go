package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
        "strings"
)

const DirectoryMimeType = "application/vnd.google-apps.folder"

type MkdirArgs struct {
	Out         io.Writer
	Name        string
	Description string
	Parents     []string
}

func (self *Drive) Mkdir(args MkdirArgs) error {
	f, err := self.mkdir(args)
	if err != nil {
		return err
	}
	fmt.Fprintf(args.Out, "Directory %s created\n", f.Id)
	return nil
}

func (self *Drive) mkdir(args MkdirArgs) (*drive.File, error) {
        // Checking for pre-existing directory
        parentId := args.Parents[0]
        name := strings.Replace(args.Name, "'", "\\'", -1)
	query := fmt.Sprintf("trashed = false and mimeType = '%s' and name = '%s' and '%s' in parents", DirectoryMimeType, name, parentId)
	result, err := self.service.Files.List().Q(query).Fields("files(id,name)").Do()
	if err != nil {
		return nil, err
	}
	if len(result.Files) == 0 {
                // Directory does not exist yet, create it
		return self.mkdirActual(args)
	}
	// Matching directory already exist, using it instead
        return result.Files[0], nil
}

func (self *Drive) mkdirActual(args MkdirArgs) (*drive.File, error) {
	dstFile := &drive.File{
		Name:        args.Name,
		Description: args.Description,
		MimeType:    DirectoryMimeType,
	}

	// Set parent folders
	dstFile.Parents = args.Parents

	// Create directory
	f, err := self.service.Files.Create(dstFile).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to create directory: %s", err)
	}

	return f, nil
}
