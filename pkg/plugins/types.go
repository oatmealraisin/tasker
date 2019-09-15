// Tasker - A pluggable task server for keeping track of all those To-Do's
// Copyright (C) 2019 Ryan Murphy <ryan@oatmealrais.in>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package plugins

import (
	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/oatmealraisin/tasker/pkg/storage"
	"github.com/spf13/cobra"
)

/* TaskerPlugin is the base interface that all plugins have to implement */
type TaskerPlugin interface {
	/* Initialize is called whenever the plugin is loaded, normally at the
	   beginning of a command. */
	Initialize() error
	/* Destroy is called when tasker stops running. Perform all cleanup needs. */
	Destroy() error
	/* Install is called once when the `tasker plugin install` command is run,
	after the plugin is built and autoloaded. */
	Install() error
	/* Uninstall is called once when the `tasker plugin uninstall` */
	Uninstall() error

	/* Name is used for formatting, you should return the pretty-name of your
	plugin. */
	Name() string
	/* Description is used for formatting, you should return a long-form
	description of what your plugin does. */
	Description() string
	/* Help is used for formatting, you should return configuration and command
	information about your plugin */
	Help() string
	/* Version is used for comparing installed plugins. */
	Version() string
}

/* A TaskModifierPlugin has the ability to influnce a task as it is processed.
 */
type TaskModifier interface {
	/* TaskCreatedHook is called when a new task is created. Specifically, it is
	called after all the internal checks and processing is done, and just before
	it is inserted into the database. Input is a copy of the task. The hook
	should return a new task with any modifications the Plugin may want to make.
	These changes will then be checked for legality. */
	TaskCreatedHook(task models.Task, get storage.GetFunc) (models.Task, error)
	/* TaskFinishedHook is called when a task is finished (removed or not).
	Specifically, it is called after all modifications have been made by Tasker,
	but before the changes are submitted to the database. Input is a copy of the
	task. The hook should return a new task with any modifications the Plugin
	may want to make. These changes will then be checked for legality. */
	TaskFinishedHook(task models.Task, get storage.GetFunc) (models.Task, error)
	/* TaskModifiedHook is called when a task is changed by the user.
	Specifically, it is called after all modifications have been made by Tasker,
	but before the changes are submitted to the database. Input is a copy of the
	task. The hook should return a new task with any modifications the Plugin
	may want to make. These changes will then be checked for legality. */
	TaskModifiedHook(task models.Task, get storage.GetFunc) (models.Task, error)
	/* ScoreModifierHook is called when the score is calculated for a task.
	Specifically, each plugin is called in parallel, and the score changes are
	added together to the final score. */
	ScoreModifierHook(task models.Task, get storage.GetFunc) (models.Task, error)
}

type TaskCreator interface {
	SetCreateFunc(add storage.CreateFunc)
}

type TaskEditor interface {
	SetEditFunc(edit storage.EditFunc)
}

type TaskViewer interface {
	SetGetFunc(get storage.GetFunc)
}

/* A TaskManager plugin has the ability to arbitrarily add, remove, and modify
Tasks async of the main Tasker objective. */
type TaskManager interface {
	TaskModifier
	TaskEditor
	TaskCreator
	TaskViewer
}

type TaskerCommand interface {
	Commands() []*cobra.Command
}
