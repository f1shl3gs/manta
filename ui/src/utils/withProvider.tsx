import React from 'react'

const withProvider = (Provider: React.FC, Component: React.FC): React.FC => {
  return () => (
    <Provider>
      <Component />
    </Provider>
  )
}

export default withProvider
