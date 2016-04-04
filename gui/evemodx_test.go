package main

import (
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
    "math/rand"
    "time"
    "log"
    "fmt"
    "io/ioutil"
    emx "evemodx"

)

var (
    logView *walk.TextEdit
    mainWindow *walk.MainWindow
    btnInject, btnRefreshPids, btnInjectAll *walk.PushButton
    modViewTable *walk.TableView
    pidViewTable *walk.TableView 
    modLabel *walk.Label
    targetMods []string
    targetPids []int
    mods []ModInfo
    pids []int
)

const (
    VERSION = "0.0.2"
)

type ModView struct {
    Index int
    ModName string
    ModifiedTime time.Time
    Description string
    checked bool
}

type ModViewModel struct {
    walk.SortedReflectTableModelBase // Reflection
    items []*ModView
}

type PidView struct {
    Index int
    Pid int
    ProcessName string
    CharName string
    checked bool
}

type PidViewModel struct {
    walk.SortedReflectTableModelBase // Reflection
    items []*PidView
}

type ModInfo struct {
    ModName string
    ModifiedTime time.Time
}

func printLog(info string){
    logView.AppendText(fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), info)  + "\r\n" )
}

func GetMods() []ModInfo {
    modReaderDir, _ := ioutil.ReadDir("./mods/")
    var mods []ModInfo
    for _, fileInfo := range modReaderDir {
        if fileInfo.IsDir() {
            mods = append(mods, ModInfo{ModName: fileInfo.Name(), ModifiedTime: fileInfo.ModTime()})
        }
    }
    return mods
}



func (n *PidViewModel) Items() interface{} {
    return n.items
}

func NewPidViewModel() *PidViewModel {
    n := new(PidViewModel)
    n.LoadPids()
    return n
}

func (n *PidViewModel) LoadPids() {

    pids = emx.GetGamePids()
    fmt.Println("%s", pids)
    n.items = make([]*PidView, len(pids))
    for i := range n.items {
        n.items[i] = &PidView{
            Index: i,
            Pid: pids[i],
            ProcessName:  "exefile.exe",
            CharName: "ISD 双鱼座",

        }
        fmt.Println("%d", pids[i])
    }
    n.PublishRowsReset()
}

func (n *PidViewModel) Checked(row int) bool {
    return n.items[row].checked
}

func (n *PidViewModel) SetChecked(row int, checked bool) error {
    n.items[row].checked = checked
    fmt.Println(fmt.Sprintf("%d %s", row, checked))
    if checked {
        targetPids = append(targetPids, n.items[row].Pid)
    } else {
        for ix, value := range targetPids {
            if value == n.items[row].Pid {
                ixa := ix + 1
                targetPids = append(targetPids[:ix], targetPids[ixa:]...)
            }
        }
    }
    fmt.Println(fmt.Sprintf("%s", targetPids))
    return nil
}






func (m *ModViewModel) Items() interface{} {
    return m.items
}

func NewModViewModel() *ModViewModel {
    m := new(ModViewModel)
    m.LoadMods()
    return m
}

func (m *ModViewModel) LoadMods() {
    mods = GetMods()
    //now := time.Now()
    m.items = make([]*ModView, len(mods))
    //fmt.Println(mods[0])
    for i := range m.items {
        m.items[i] = &ModView{
            Index: i,
            ModName:  mods[i].ModName,
            ModifiedTime:   mods[i].ModifiedTime,
            Description: "This is a mod.",

        }
    }
    m.PublishRowsReset()
}

func (m *ModViewModel) Checked(row int) bool {
    return m.items[row].checked
}

func (m *ModViewModel) SetChecked(row int, checked bool) error {
    m.items[row].checked = checked
    fmt.Println(fmt.Sprintf("%d %s", row, checked))
    if checked {
        targetMods = append(targetMods, m.items[row].ModName)
    } else {
        for ix, value := range targetMods {
            if value == m.items[row].ModName {
                ixa := ix + 1
                targetMods = append(targetMods[:ix], targetMods[ixa:]...)
            }
        }
    }
    fmt.Println(fmt.Sprintf("%s", targetMods))
    return nil
}



func main() {


    rand.Seed(time.Now().UnixNano())
    modViewModel := NewModViewModel()
    pidViewModel := NewPidViewModel()
    currentModDirectory := emx.GetCurrentDirectory() + "/mods/"

    if err := (MainWindow{
        AssignTo: &mainWindow,
        Title:   "EVEModX",
        MaxSize: Size{1080, 600},
        MinSize: Size{1080, 600},

        Layout:  Grid{},
        Children: []Widget{

            Label{
                Row: 0,
                Column: 0,
                AssignTo: &modLabel,
                Text: "Mods",
            },

            HSpacer{Row: 0, Column: 1, ColumnSpan: 15},

            TableView{
                Row: 1,
                Column: 0,
                ColumnSpan: 16,
                AssignTo: &modViewTable,
                CheckBoxes:            true,
                ColumnsOrderable:      false,
                MultiSelection:        true,
                Columns: []TableViewColumn{
                    {DataMember: "Index", Width: 45},
                    {DataMember: "ModName", Width: 170},
                    {DataMember: "ModifiedTime", Format: "2006-01-02 15:04:05", Width: 150},
                    {DataMember: "Description", Width: 650},
                },
                Model: modViewModel,
                OnCurrentIndexChanged: func() {
                    
                },
            },


            TableView{
                Row: 2,
                Column: 0,
                ColumnSpan: 6,
                AssignTo: &pidViewTable,
                CheckBoxes:            true,
                ColumnsOrderable:      false,
                MultiSelection:        true,
                Columns: []TableViewColumn{
                    {DataMember: "Index"},
                    {DataMember: "Pid"},
                    {DataMember: "ProcessName"},
                    {DataMember: "CharName"},
                },
                Model: pidViewModel,
                OnCurrentIndexChanged: func() {
                    
                },
            },


            TextEdit{
                AssignTo: &logView,
                Row: 2,
                Column: 6,
                ColumnSpan: 10,
                ReadOnly: true,
            },


            PushButton{
                Row: 3,
                Column: 0,
                
                AssignTo:  &btnInject,
                Text: "Execute",
                MinSize: Size{76, 26},
                MaxSize: Size{76, 26},
                OnClicked: func() {
                    var code string
                    for _, value := range targetMods {
                        code = code + "import " + value + ";"
                    }
                    code = `import sys;sys.path.append('` + currentModDirectory + `');` + code + ``
                    fmt.Println(fmt.Sprintf("%s", code))
                    for i, _ := range targetPids {
                        pid := fmt.Sprintf("%d",targetPids[i])
                        printLog(fmt.Sprintf("[INFO] Executing injection for %d", targetPids[i]))
                        emx.Inject(pid, code)

                    }

                },

            },  
            PushButton{
                Row: 3,
                Column: 1,
                
                AssignTo:  &btnInjectAll,
                Text: "Execute All PIDs",
                MinSize: Size{76, 26},
                MaxSize: Size{76, 26},
                OnClicked: func() {
                    /*var code string
                    for _, value := range mods {
                        code = code + "import " + value.ModName + ";"
                    }*/
                    var code string
                    for _, value := range targetMods {
                        code = code + "import " + value + ";"
                    }                    
                    code = `import sys;sys.path.append('` + currentModDirectory + `');` + code + ``
                    fmt.Println(fmt.Sprintf("%s", code))
                    for i, _ := range pids {
                        pid := fmt.Sprintf("%d",pids[i])
                        printLog(fmt.Sprintf("[INFO] Executing injection for %d", pids[i]))
                        emx.Inject(pid, code)

                    }
                },

            },                   
            PushButton{
                Row: 3,
                Column: 2,
                //ColumnSpan: 2,
                AssignTo:  &btnRefreshPids,
                Text: "Reload PID",
                MinSize: Size{76, 26},
                MaxSize: Size{76, 26},
                OnClicked: pidViewModel.LoadPids,

            },         

        },
    }).Create(); err != nil {
        log.Fatal(err)
    }


    modLabel.SetSize(walk.Size{50, 18})

    printLog("[INFO] EVEModX " + VERSION + " started")
    printLog("[INFO] init")
    mainWindow.Run()
}


