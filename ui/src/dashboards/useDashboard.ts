import constate from 'constate'
import {Dashboard} from 'src/types/Dashboard'

interface State {
  dashboard: Dashboard
}

const [DashboardProvider, useDashboard] = constate((state: State) => {
  return {
    ...state.dashboard,
  }
})

export {DashboardProvider, useDashboard}
