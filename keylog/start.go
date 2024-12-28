package keylog

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/keylogme/zero-trust-logger/keylog/storage"
)

func Start(
	chEvt chan DeviceEvent,
	devices *[]Device,
	sd *shortcutsDetector,
	store storage.Storage,
) {
	slog.Info("Listening...")

	go func() {
		for i := range chEvt {
			sd := sd.handleKeyEvent(i)
			if sd.ShortcutId != 0 {
				slog.Info(
					fmt.Sprintf(
						"Shortcut %d found in device %s\n",
						sd.ShortcutId,
						sd.DeviceId,
					),
				)
				store.SaveShortcut(sd.DeviceId, sd.ShortcutId)
			}

			if i.Type == evKey && i.KeyRelease() {
				start := time.Now()
				// FIXME: mod+key is sent, but when mod is released , is sent again
				// keylogs := []uint16{i.Code}
				// keylogs = append(keylogs, modPress...)
				// err := sendKeylog(sender, i.DeviceId, i.Code)
				err := store.SaveKeylog(i.DeviceId, i.Code)
				if err != nil {
					fmt.Printf("error %s\n", err.Error())
				}
				slog.Info(
					fmt.Sprintf(
						"| %s | Key :%d %s\n",
						time.Since(start),
						i.Code,
						i.KeyString(),
					),
				)
			}
		}
	}()
}
