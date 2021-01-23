import {Dispatch, SetStateAction, useCallback, useState} from 'react'

const useLocalStorage = <T>(
  key: string,
  initialValue: T,
  raw?: boolean
): [T | undefined, Dispatch<SetStateAction<T | undefined>>] => {
  // eslint-disable-next-line react-hooks/exhaustive-deps
  let deserializer = raw && !raw ? (value: string) => value : JSON.parse
  // eslint-disable-next-line react-hooks/exhaustive-deps
  let serializer = raw && !raw ? (value: string) => value : JSON.stringify

  const [stored, setStored] = useState<T | undefined>(() => {
    try {
      const storedValue = localStorage.getItem(key)
      if (storedValue !== null) {
        return deserializer(storedValue)
      }

      initialValue && localStorage.setItem(key, serializer(initialValue))
      return initialValue
    } catch (error) {
      return initialValue
    }
  })

  const set: Dispatch<SetStateAction<T | undefined>> = useCallback(
    (valOrFunc) => {
      try {
        const next =
          typeof valOrFunc === 'function'
            ? (valOrFunc as Function)(stored)
            : valOrFunc
        if (typeof next === 'undefined') {
          return
        }

        localStorage.setItem(key, serializer(next))
        setStored(deserializer(next))
      } catch (error) {
        console.error(error)
      }
    },
    [key, stored, serializer, deserializer]
  )

  return [stored, set]
}

export default useLocalStorage
