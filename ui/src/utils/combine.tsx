import React from 'react'

type ProviderWithValue<T = any> =
  | React.ComponentType<T>
  | [React.ComponentType<T>, T]

/*
* Example:
const CombinedProviders = combineProviders([
  [Provider1, { initialState: 5 }],
  Provider2,
]);
*
* */
const combineProviders = (providers: Array<ProviderWithValue>) => {
  return ({children}: React.PropsWithChildren<{value?: any[]}>) => {
    return providers.reduce<React.ReactElement<React.ProviderProps<any>>>(
      (tree, PV) => {
        if (Array.isArray(PV)) {
          const [Provider, value] = PV
          return <Provider {...value}>{tree}</Provider>
        }

        return <PV>{tree}</PV>
      },
      children as React.ReactElement
    )
  }
}

export default combineProviders
