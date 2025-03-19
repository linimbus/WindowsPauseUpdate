package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var db *walk.DataBinder
var mainWindow *walk.MainWindow
var active *walk.PushButton
var pauseWeek *walk.NumberEdit
var start bool
var featureStart, featureEnd, qualityStart, qualityEnd *walk.Label

func init() {
	go AutoStartup()
	go ActiveTask()
}

func AutoStartup() {
	for {
		if mainWindow != nil && mainWindow.Visible() {
			break
		}
		time.Sleep(time.Second)
	}

	ActiveEnable()
}

func ActiveTask() {
	for {
		time.Sleep(time.Hour)
		if start {
			ActiveOnce()
		}
	}
}

func ActiveOnce() {

	weeks := ConfigGet().PauseWeeks
	start, end := time.Now(), time.Now()
	start = start.Add(-time.Hour * 24)
	days := 7*time.Duration(weeks) - 1
	end = end.Add(time.Hour * 24 * days)

	if configCache.QualityEndTime.Compare(end) > 0 {
		StatusUpdate(configCache.QualityEndTime.Format(time.DateTime))
		return
	}

	configCache.FeatureEndTime = end
	configCache.FeatureStartTime = start
	configCache.QualityEndTime = end
	configCache.QualityStartTime = start

	err := RegistryWrite(configCache)
	if err != nil {
		StatusUpdate(err.Error())
		return
	}

	StatusUpdate(configCache.QualityEndTime.Format(time.DateTime))

	featureStart.SetText(configCache.FeatureStartTime.Format(time.DateTime))
	featureEnd.SetText(configCache.FeatureEndTime.Format(time.DateTime))
	qualityStart.SetText(configCache.QualityStartTime.Format(time.DateTime))
	qualityEnd.SetText(configCache.FeatureEndTime.Format(time.DateTime))
}

func ActiveEnable() {
	active.SetEnabled(false)
	defer active.SetEnabled(true)

	start = !start

	logs.Info("config %s", ConfigString())

	if start {
		ActiveOnce()
	}

	if start {
		active.SetImage(ICON_Stop)
		active.SetText("Stop")
	} else {
		active.SetImage(ICON_Start)
		active.SetText("Start")
	}

	pauseWeek.SetEnabled(!start)
}

func MainWindows() {
	CapSignal(CloseWindows)

	err := ConfigInit()
	if err != nil {
		logs.Error("config init failed, %s", err.Error())
		StatusUpdate(err.Error())
	}

	config := ConfigGet()

	cnt, err := MainWindow{
		Title:    "Windows Pause Update " + VersionGet(),
		Icon:     ICON_Main,
		AssignTo: &mainWindow,
		MinSize:  Size{Width: 400, Height: 250},
		Size:     Size{Width: 400, Height: 250},
		Layout:   Grid{Columns: 2, Margins: Margins{Top: 5, Bottom: 10, Left: 10, Right: 10}},
		DataBinder: DataBinder{
			AssignTo:   &db,
			Name:       "",
			DataSource: config,
		},
		MenuItems: []MenuItem{
			Action{
				Text: "Runlog",
				OnTriggered: func() {
					OpenBrowserWeb(RunlogDirGet())
				},
			},
			Action{
				Text: "About",
				OnTriggered: func() {
					AboutAction()
				},
			},
		},
		StatusBarItems: StatusBarInit(),
		Children: []Widget{
			Label{
				Text: "Feature Update Start Time: ",
			},
			Label{
				AssignTo: &featureStart,
				Text:     config.FeatureStartTime.Format(time.DateTime),
			},

			Label{
				Text: "Feature Update End Time: ",
			},
			Label{
				AssignTo: &featureEnd,
				Text:     config.FeatureEndTime.Format(time.DateTime),
			},

			Label{
				Text: "Quality Update Start Time: ",
			},
			Label{
				AssignTo: &qualityStart,
				Text:     config.QualityStartTime.Format(time.DateTime),
			},

			Label{
				Text: "Quality Update End Time: ",
			},
			Label{
				AssignTo: &qualityEnd,
				Text:     config.QualityEndTime.Format(time.DateTime),
			},

			Label{
				Text: "Pause Week: ",
			},
			NumberEdit{
				AssignTo:  &pauseWeek,
				Value:     config.PauseWeeks,
				Suffix:    " Weeks",
				Decimals:  0,
				Increment: 1,
				MinValue:  1.0,
				MaxValue:  5.0,
			},

			PushButton{
				AssignTo:   &active,
				Text:       "Start",
				Image:      ICON_Start,
				ColumnSpan: 2,
				OnClicked: func() {
					ActiveEnable()
				},
			},
		},
	}.Run()

	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("main windows exit %d", cnt)
	}

	if err := recover(); err != nil {
		logs.Error(err)
	}

	CloseWindows()
}

func CloseWindows() {
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
}
