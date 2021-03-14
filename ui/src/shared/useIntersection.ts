import {RefObject, useEffect, useState} from 'react'
import {useSearch} from './useSearch'

/*
Example:
const App = () => {
  // Ref for the element that we want to detect whether on screen
  const ref = useRef();

  // Call the hook passing in ref and root margin
  // In this case it would only be considered onScreen if more ...
  // ... than 300px of element is visible.
  const intersecting = useIntersecting(ref, '-300px');

  return (
    <div ref={ref}>
      {intersecting ? <p>on screen</p> : <p>not</p>}
    </div>
  )
}
* */
function useIntersection(
  ref: RefObject<HTMLElement>,
  rootMargin = '0px'
): boolean {
  // State and setter for storing whether element is visible
  const [intersecting, setIntersecting] = useState(false)

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        // Update our state when observer callback fires
        setIntersecting(entry.isIntersecting)
      },
      {
        rootMargin,
      }
    )

    if (ref.current) {
      observer.observe(ref.current)
    }

    return () => {
      observer.disconnect()
    }
    // Empty array ensures that effect is only run on mount and unmount
  }, [rootMargin])

  return intersecting
}

export default useIntersection
