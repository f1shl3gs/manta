import {ResourceType} from 'src/types/resources'
import {FunctionComponent, ReactNode} from 'react'
import PageSpinner from 'src/shared/components/PageSpinner'
import {AppState} from 'src/types/stores'
import {getResourcesStatus} from '../selectors';

interface Props {
  resources: Array<ResourceType>
  childrent: ReactNode
}

const GetResources: FunctionComponent<Props> = ({resources, children}) => {

  return (
    <PageSpinner loading={loading}>
      {children}
    </PageSpinner>
  )
}

const mstp = (state: AppState, {resources}: Props) => {
  const loading = getResourcesStatus(state, resources)

  return {
    loading
  }
}

const mdtp = {
  getDashboards: getDashboards
}

const connector = connect(mstp, mdtp)

export default GetResources
