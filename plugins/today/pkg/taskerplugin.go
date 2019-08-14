// Tasker - A pluggable task server for keeping track of all those To-Do's
// Today - A plugin for focusing on a subset of tasks just for today
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
package today

func (t *Today) Initialize() error {
	//now := time.Now()
	//now_s := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	//if _, ok := t.tasks[now_s]; !ok {
	//	fmt.Println("Nothing to do today")
	//}

	//t.today = t.tasks[now_s]

	t.initialized = true
	return nil
}

func (t *Today) Destroy() error {
	return nil
}

func (t *Today) Install() error {
	return nil
}

func (t *Today) Uninstall() error {
	return nil
}

func (t *Today) Name() string {
	return "Today"
}

func (t *Today) Description() string {
	panic("not implemented")
}

func (t *Today) Help() string {
	panic("not implemented")
}

func (t *Today) Version() string {
	return "0.1.0"
}
