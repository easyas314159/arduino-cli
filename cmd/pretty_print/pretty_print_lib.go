/*
 * This file is part of arduino-cli.
 *
 * arduino-cli is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2017 BCMI LABS SA (http://www.arduino.cc/)
 */

package prettyPrints

import (
	"fmt"

	"github.com/bcmi-labs/arduino-cli/cmd/formatter"
	"github.com/bcmi-labs/arduino-cli/libraries"
	"github.com/bcmi-labs/arduino-cli/task"
)

// LibStatus pretty prints libraries from index status.
func LibStatus(status *libraries.StatusContext, verbosity int) {
	message := ""

	for _, name := range status.Names() {
		if verbosity > 0 {
			lib := status.Libraries[name]
			message += fmt.Sprint(lib)
			if verbosity > 1 {
				for _, r := range lib.Releases {
					message += fmt.Sprint(r)
				}
			}
			message += "\n"
		} else {
			message += name
		}
	}
	formatter.Print(formatter.Message{
		Header: "Library Search Results:",
		Data:   message,
	})
}

// DownloadLibFileIndex shows info regarding the download of a missing (or corrupted) file index of libraries.
func DownloadLibFileIndex() task.Wrapper {
	return DownloadFileIndex(libraries.DownloadLibrariesFile)
}

//UninstallLib pretty prints info about a pending install of libraries.
func UninstallLib(libraryOK []string, libraryFails map[string]string) {
	actionOnItems("libraries", "uninstalled", libraryOK, libraryFails)
}

// CorruptedLibIndexFix pretty prints messages regarding corrupted index fixes of libraries.
func CorruptedLibIndexFix(index libraries.Index, verbosity int) (libraries.StatusContext, error) {
	downloadTask := DownloadLibFileIndex()
	parseTask := libIndexParse(index, verbosity)

	result := corruptedIndexFixResults(downloadTask, parseTask, verbosity)

	return result[1].Result.(libraries.StatusContext), result[0].Error
}

// libIndexParse pretty prints info about parsing an index file of libraries.
func libIndexParse(index libraries.Index, verbosity int) task.Wrapper {
	ret := indexParseWrapperSkeleton()
	ret.Task = task.Task(func() task.Result {
		index.CreateStatusContext()
		return task.Result{
			Result: nil,
			Error:  nil,
		}
	})
	return ret
}
