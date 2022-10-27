// Libraries
import constate from 'constate'
import {useState} from 'react'

const [StepProvider, useStep] = constate(() => {
  const [step, setStep] = useState(0)
  const next = () => setStep(prevState => prevState + 1)

  return {
    step,
    next,
  }
})

export {StepProvider, useStep}
