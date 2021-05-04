import React from 'react'

const withProvider = (
  Provider: React.FC,
  Component: React.FC
): React.ComponentType<any> => {
  return () => (
    <Provider>
      <Component />
    </Provider>
  )
}

export default withProvider
