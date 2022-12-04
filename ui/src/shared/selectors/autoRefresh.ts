import {useSelector} from 'react-redux'
import {AppState} from 'src/types/stores'

export function useAutoRefresh() {
  useSelector((state: AppState) => {
    const {start, end, step} = state.autoRefresh

    return {
      start,
      end,
      step,
    }
  })
}
