import {throttle} from 'lodash'
import {Store} from 'redux'

import {LocalStorage} from 'src/types/localStorage'
import { saveToLocalStorage } from 'src/store/localStorage'

export default function persistState() {
  return next => (reducer, initialState: LocalStorage, enhancer) => {
    const store: Store<LocalStorage> = next(reducer, initialState, enhancer)
    const throttleMs = 1000

    store.subscribe(
      throttle(() => {
        saveToLocalStorage(store.getState())
      }, throttleMs)
    )

    return store
  }
}
