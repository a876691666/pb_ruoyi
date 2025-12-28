package system

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func getSuffix(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i+1:]
		}
	}
	return ""
}

func syncOSS(e *core.RecordRequestEvent) error {
	files := e.Record.GetUnsavedFiles("file")

	if len(files) == 0 {
		return e.Next()
	}

	file := files[0]

	if file == nil {
		return e.Next()
	}

	e.Record.Set("file_name", file.Name)
	e.Record.Set("original_name", file.OriginalName)
	e.Record.Set("file_suffix", getSuffix(file.OriginalName))

	return e.Next()
}

func syncOSSAfter(e *core.RecordEvent) error {

	url := "/api/files/" + e.Record.BaseFilesPath() + "/" + e.Record.GetString("file")
	if e.Record.GetString("url") != url {
		e.Record.Set("url", url)
		if err := e.App.Save(e.Record); err != nil {
			return err
		}
	}

	return e.Next()
}

// RegisterSystemOSS 注册
func RegisterSystemOSS(app *pocketbase.PocketBase) {
	app.OnRecordCreateRequest("oss").BindFunc(syncOSS)
	app.OnRecordUpdateRequest("oss").BindFunc(syncOSS)

	app.OnRecordAfterCreateSuccess("oss").BindFunc(syncOSSAfter)
}
