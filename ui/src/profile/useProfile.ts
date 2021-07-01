import constate from 'constate'

import testData from './components/testData'

const [ProfileProvider, useProfile] = constate(
  () => {
    return testData
  },
  values => values
)

export {ProfileProvider, useProfile}
