import React, {useCallback} from 'react'
import {Route, Switch, useHistory, useParams} from 'react-router-dom'

import {
  Button,
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
} from '@influxdata/clockface'
import Otcls from './otcls'
import OtclEdit from './component/OtclEdit'
import OtclCreate from './component/OtclCreate'

import {OtclProvider} from './state'
import {useOrgID} from '../shared/state/organization/organization'

const pageContentsClassName = `alerting-index alerting-index__${'check'}`
const title = 'OpenTelemetry Collector'
const otclsPrefix = `/orgs/:orgID/otcls`

const Header: React.FC = () => {
  return (
    <Page.Header fullWidth>
      <Page.Title title={title} />
    </Page.Header>
  )
}

type OtclPageProps = {
  onCreate: () => void
}

class OtclPage extends React.PureComponent<OtclPageProps> {
  render() {
    const {onCreate} = this.props

    return (
      <Page titleTag={title}>
        <Header />
        <Page.ControlBar fullWidth>
          <Page.ControlBarRight>
            <Button
              size={ComponentSize.Small}
              icon={IconFont.Plus}
              color={ComponentColor.Primary}
              text="Create Configuration"
              onClick={onCreate}
            />
          </Page.ControlBarRight>
        </Page.ControlBar>

        <Page.Contents
          fullWidth
          scrollable={false}
          className={pageContentsClassName}
        >
          <Otcls />
        </Page.Contents>
      </Page>
    )
  }
}

const Otcl: React.FC = () => {
  const orgID = useOrgID()
  const history = useHistory()
  const onCreate = useCallback(() => {
    history.push(`/orgs/${orgID}/otcls/new`)
  }, [orgID])

  return (
    <OtclProvider orgID={orgID}>
      <OtclPage onCreate={onCreate} />
      <Switch>
        <Route path={`${otclsPrefix}/new`} component={OtclCreate} />
        <Route path={`${otclsPrefix}/:otclID`} component={OtclEdit} />
      </Switch>
    </OtclProvider>
  )
}

export default Otcl
