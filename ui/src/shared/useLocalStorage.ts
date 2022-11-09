import {useState} from 'react'

function useLocalStorage<T>(key: string, initialValue: T) {
  const [storedValue, setStoredValue] = useState<T>(() => {
    try {
      // Get from local storage by key
      const item = window.localStorage.getItem(key)
      return item ? JSON.parse(item) : initialValue
    } catch (err) {
      console.error(
        `read value from local storage failed, key: ${key}, err: ${err}`
      )
      return initialValue
    }
  })

  const setValue = (value: T | ((value: T) => T)) => {
    try {
      // Allow value to be a function so we have same API as useState
      const valueToStore =
        value instanceof Function ? value(storedValue) : value
      // save state
      setStoredValue(valueToStore)
      // save to local storage
      window.localStorage.setItem(key, JSON.stringify(valueToStore))
    } catch (err) {
      console.error(
        `set key value to local storage failed, key: ${key}, err: ${err}`
      )
    }
  }

  return [storedValue, setValue] as const
}

export default useLocalStorage
