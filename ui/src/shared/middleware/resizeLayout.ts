import {Dispatch, Action, Middleware} from 'redux'

// Trigger resize event to re-layout the React Layout plugin
export const resizeLayout: Middleware =
  () => (next: Dispatch<Action>) => (action: Action) => {
    next(action)

    if (
      action.type === 'EnablePresentationMode' ||
      action.type === 'DisablePresentationMode'
    ) {
      setTimeout(() => {
        // Uses longer event object creation method due to IE compatibility.
        const ev = document.createEvent('HTMLEvents')
        ev.initEvent('resize', false, true)
        window.dispatchEvent(ev)
      }, 50)
    }
  }
