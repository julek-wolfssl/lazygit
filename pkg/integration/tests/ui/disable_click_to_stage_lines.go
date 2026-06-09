package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableClickToStageLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With disableClickToStageLines set, clicking a diff doesn't focus it; only a key does",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.DisableClickToStageLines = true
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("file1"))

		t.Views().Main().
			Click(1, 8).
			SelectionIsHidden()

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// Once focused with a key, clicks select lines as usual.
		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+three")).
			Click(1, 8).
			SelectedLines(Contains("+four")).
			PressPrimaryAction().
			PressEscape()

		// The staged changes are in the secondary pane, which doesn't take the focus
		// on a click either.
		t.Views().Secondary().
			Content(Contains("+four")).
			Click(1, 7).
			SelectionIsHidden()

		t.Views().Files().
			IsFocused()

		// A main view showing no diff still takes the focus on a click.
		t.Views().Branches().
			Focus()

		t.Views().Main().
			Click(1, 0).
			IsFocused()
	},
})
