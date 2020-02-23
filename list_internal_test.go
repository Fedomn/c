package main

import (
	"runtime"
	"time"

	keybd "github.com/micmonay/keybd_event"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// simulate keyboard events
var kb keybd.KeyBonding

var _ = BeforeSuite(func() {
	var err error
	kb, err = keybd.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	// For linux, it is very important wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
})

func pressKey(keys ...int) {
	kb.SetKeys(keys...)
	err := kb.Launching()
	if err != nil {
		panic(err)
	}
	// clear ctrl shift alt keys
	kb.Clear()
}

func pressKeyWithCtrl(keys ...int) {
	kb.HasCTRL(true)
	pressKey(keys...)
}

func pressKeyWithShift(keys ...int) {
	kb.HasSHIFT(true)
	pressKey(keys...)
}

// simulate real keyboard operation
var _ = Describe("Select List", func() {

	var (
		cmds       []Cmd
		cmdChan    chan Cmd
		selectList *SelectList
	)

	BeforeEach(func() {
		cmds = []Cmd{
			{Name: "normal_cmd1_name", Cmd: `echo normal_cmd1_name`},
			{Name: "normal_cmd2_name", Cmd: `echo normal_cmd2_name`},
			{Name: "search_cmd1_name", Cmd: `echo search_cmd1_name`},
			{Name: "search_cmd2_name", Cmd: `echo search_cmd2_name`},
		}
		cmdChan = make(chan Cmd)
		selectList = NewUIList(cmds, cmdChan)
	})

	AfterEach(func() {
		selectList.close()
	})

	Context("Normal Mode", func() {
		It("should initUI with expected arguments", func() {
			Expect(selectList.uiList.Title).To(HavePrefix("Usage:"))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			Expect(selectList.uiList.Rows).To(HaveLen(4))
			Expect(selectList.uiList.Rows[0]).To(ContainSubstring("normal_cmd1_name"))
			Expect(selectList.uiList.Rows[2]).To(ContainSubstring("search_cmd1_name"))
		})

		It("should scroll down by shortcut <j>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKey(keybd.VK_J)
			Expect(selectList.uiList.SelectedRow).To(Equal(1))

			close(done)
		})

		It("should scroll up by shortcut <k>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKey(keybd.VK_J)
			Expect(selectList.uiList.SelectedRow).To(Equal(1))
			pressKey(keybd.VK_K)
			Expect(selectList.uiList.SelectedRow).To(Equal(0))

			close(done)
		})

		It("should scroll page down by shortcut <C-f>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKeyWithCtrl(keybd.VK_F)
			Expect(selectList.uiList.SelectedRow).To(Equal(3))

			close(done)
		})

		It("should scroll page up by shortcut <C-b>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKeyWithCtrl(keybd.VK_F)
			Expect(selectList.uiList.SelectedRow).To(Equal(3))
			pressKeyWithCtrl(keybd.VK_B)
			Expect(selectList.uiList.SelectedRow).To(Equal(0))

			close(done)
		})

		It("should close UI by shortcut <q>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKey(keybd.VK_J)
			Expect(selectList.uiList.SelectedRow).To(Equal(1))
			pressKey(keybd.VK_Q)

			Expect(<-cmdChan).To(Equal(Cmd{}))
			Expect(selectList.isClose).To(BeTrue())

			close(done)
		})

		It("should send cmd to chan by shortcut <Enter>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKey(keybd.VK_J)
			Expect(selectList.uiList.SelectedRow).To(Equal(1))
			pressKey(keybd.VK_ENTER)

			Expect(<-cmdChan).To(Equal(Cmd{Name: "normal_cmd2_name", Cmd: `echo normal_cmd2_name`}))
			Expect(selectList.isClose).To(BeTrue())

			close(done)
		})

		It("should into search mode by shortcut </>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKey(keybd.VK_J)
			Expect(selectList.uiList.SelectedRow).To(Equal(1))
			pressKey(keybd.VK_SP11)

			Expect(selectList.selectedMode).To(Equal(SearchMode))
			Expect(selectList.uiList.Title).To(HavePrefix("Search:"))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			Expect(selectList.uiList.Rows).To(HaveLen(4))
			Expect(selectList.uiList.Rows[0]).To(ContainSubstring("normal_cmd1_name"))
			Expect(selectList.uiList.Rows[2]).To(ContainSubstring("search_cmd1_name"))

			close(done)
		})

	})

	Context("Search Mode", func() {
		It("should scroll down by shortcut <C-j>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.selectedMode).To(Equal(NormalMode))

			pressKey(keybd.VK_J)
			//into search mode
			pressKey(keybd.VK_SP11)

			Expect(selectList.selectedMode).To(Equal(SearchMode))
			Expect(selectList.uiList.Title).To(HavePrefix("Search:"))
			Expect(selectList.uiList.Rows).To(HaveLen(4))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKeyWithCtrl(keybd.VK_J)
			Expect(selectList.uiList.SelectedRow).To(Equal(1))

			close(done)
		})

		It("should scroll up by shortcut <C-k>", func(done Done) {
			go selectList.ListenEvents()

			Expect(selectList.selectedMode).To(Equal(NormalMode))

			pressKey(keybd.VK_J)
			//into search mode
			pressKey(keybd.VK_SP11)

			Expect(selectList.selectedMode).To(Equal(SearchMode))
			Expect(selectList.uiList.Title).To(HavePrefix("Search:"))
			Expect(selectList.uiList.Rows).To(HaveLen(4))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			pressKeyWithCtrl(keybd.VK_J)
			pressKeyWithCtrl(keybd.VK_J)
			pressKeyWithCtrl(keybd.VK_K)
			Expect(selectList.uiList.SelectedRow).To(Equal(1))

			close(done)
		})

		It("should go through search flow normally and end with <Enter>", func(done Done) {
			go selectList.ListenEvents()

			pressKey(keybd.VK_J)
			//into search mode
			pressKey(keybd.VK_SP11)

			// press "search_cmd2" to search "search_cmd2_name"
			// search
			pressKey(keybd.VK_S, keybd.VK_E, keybd.VK_A, keybd.VK_R, keybd.VK_C, keybd.VK_H)
			// _
			pressKeyWithShift(keybd.VK_SP2)
			// cmd2
			pressKey(keybd.VK_C, keybd.VK_M, keybd.VK_D, keybd.VK_2)

			Expect(selectList.uiList.Title).To(ContainSubstring("search_cmd2"))
			Expect(selectList.uiList.Rows).To(HaveLen(1))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))

			// press <Backspace> to search "search_cmd"
			pressKey(keybd.VK_DELETE)
			Expect(selectList.uiList.Title).To(ContainSubstring("search_cmd"))
			Expect(selectList.uiList.Rows).To(HaveLen(2))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))

			// press <Space> will search nothing
			pressKey(keybd.VK_SPACE)
			Expect(selectList.uiList.Title).To(ContainSubstring("search_cmd "))
			Expect(selectList.uiList.Rows).To(HaveLen(0))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))

			// press <C-u> to erase search string
			pressKeyWithCtrl(keybd.VK_U)
			Expect(selectList.uiList.Title).To(HavePrefix("Search: "))
			Expect(selectList.uiList.Rows).To(HaveLen(4))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))

			// press <Enter> to select "normal_cmd1_name"
			pressKey(keybd.VK_ENTER)
			Expect(<-cmdChan).To(Equal(Cmd{Name: "normal_cmd1_name", Cmd: `echo normal_cmd1_name`}))
			Expect(selectList.isClose).To(BeTrue())

			close(done)
		}, 10)

		It("should go through search flow normally and end with <C-c>", func(done Done) {
			go selectList.ListenEvents()

			pressKey(keybd.VK_J)
			//into search mode
			pressKey(keybd.VK_SP11)

			// press "search_cmd2" to search "search_cmd2_name"
			// search
			pressKey(keybd.VK_S, keybd.VK_E, keybd.VK_A, keybd.VK_R, keybd.VK_C, keybd.VK_H)
			// _
			pressKeyWithShift(keybd.VK_SP2)
			// cmd2
			pressKey(keybd.VK_C, keybd.VK_M, keybd.VK_D, keybd.VK_2)

			Expect(selectList.uiList.Title).To(ContainSubstring("search_cmd2"))
			Expect(selectList.uiList.Rows).To(HaveLen(1))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))

			// press <C-c> to exit SearchMode and into NormalMode
			pressKeyWithCtrl(keybd.VK_C)
			Expect(selectList.selectedMode).To(Equal(NormalMode))
			Expect(selectList.uiList.Title).To(HavePrefix("Usage:"))
			Expect(selectList.uiList.SelectedRow).To(Equal(0))
			Expect(selectList.uiList.Rows).To(HaveLen(4))

			close(done)
		}, 10)
	})

})
