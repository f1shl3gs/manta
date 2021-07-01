import constate from 'constate'
import {useCallback, useState} from 'react'
import {useHistory, useLocation} from 'react-router-dom'

export enum ViewType {
  FlameGraph = 'FlameGraph',
  Table = 'Table',
  Both = 'Both',
}

const [ViewTypeProvider, useViewType] = constate(
  () => {
    const [viewType, setViewType] = useState(() => {
      const params = new URLSearchParams(window.location.search)
      switch (params.get('viewType')) {
        case ViewType.Both:
          return ViewType.Both
        case ViewType.Table:
          return ViewType.Table
        case ViewType.FlameGraph:
          return ViewType.FlameGraph
        default:
          return ViewType.Both
      }
    })
    const location = useLocation()
    const history = useHistory()

    const svt = useCallback(
      (vt: ViewType) => {
        setViewType(vt)
        const params = new URLSearchParams(location.search)
        params.set('viewType', vt)
        history.push(`${location.pathname}?${params.toString()}`)
      },
      [location, history]
    )

    return {
      viewType,
      setViewType: svt,
    }
  },
  values => values
)

export {ViewTypeProvider, useViewType}
