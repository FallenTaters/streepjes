package authroutine

import (
	"sync"
	"time"

	"github.com/FallenTaters/streepjes/api"
	"github.com/FallenTaters/streepjes/domain/authdomain"
	"github.com/FallenTaters/streepjes/frontend/backend"
	"github.com/FallenTaters/streepjes/frontend/backend/cache"
	"github.com/FallenTaters/streepjes/frontend/events"
	"github.com/FallenTaters/streepjes/frontend/global"
	"github.com/FallenTaters/streepjes/frontend/jscall/document"
	"github.com/FallenTaters/streepjes/frontend/jscall/window"
	"github.com/FallenTaters/streepjes/frontend/store"
)

var tracker struct {
	sync.Mutex

	lastActivity   time.Time
	lastActiveCall time.Time
}

var checkInterval time.Duration = time.Second

func Start() {
	events.Listen(events.Unauthorized, `logout-invalidate`, func() {
		document.DeleteCookie(api.AuthTokenCookieName)
		cache.InvalidateAll()
	})

	tracker.lastActivity = time.Now()

	window.OnClick(onActivity)

	go routine()
	go postActive()
}

func onActivity() {
	go check()

	tracker.Lock()
	defer tracker.Unlock()
	tracker.lastActivity = time.Now()
}

func routine() {
	for {
		check()

		time.Sleep(checkInterval)
	}
}

func check() {
	if !store.Auth.LoggedIn {
		return
	}

	tracker.Lock()
	defer tracker.Unlock()

	// logout if too long
	if time.Since(tracker.lastActivity) > authdomain.TokenDuration {
		go backend.PostLogout()
		events.Trigger(events.Unauthorized)
		return
	}

	// show warning if past warning time
	if time.Since(tracker.lastActivity) > authdomain.LockScreenWarningTime &&
		time.Since(tracker.lastActivity)-authdomain.LockScreenWarningTime < checkInterval {
		events.Trigger(events.InactiveWarning)
		return
	}

	// check if can do active call
	if tracker.lastActivity.After(tracker.lastActiveCall) {
		tracker.lastActiveCall = time.Now()

		go postActive()
	}
}

func postActive() {
	originalStore := store.Auth
	user, err := backend.PostActive()
	if err == nil {
		store.Auth.LogIn(user)
	}

	if originalStore != store.Auth {
		go global.LockAndRender()()
		events.Trigger(events.Login)
	}
}
