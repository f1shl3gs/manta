import constate from 'constate'
import {useState} from 'react'

export enum ViewType {
  FlameGraph = 'FlameGraph',
  Table = 'Table',
  Both = 'Both',
}

const [ProfileProvider, useProfile, useViewType] = constate(
  () => {
    const [viewType, setViewType] = useState(ViewType.Both)
    return {viewType, setViewType}
  },
  values => values,
  values => {
    return {
      viewType: values.viewType,
      setViewType: values.setViewType,
    }
  }
)

export {ProfileProvider, useProfile, useViewType}
