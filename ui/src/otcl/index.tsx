import React, { useCallback } from 'react';
import { Route, Switch, useHistory, useParams } from 'react-router';

import {
  Button,
  ComponentColor,
  ComponentSize,
  IconFont,
  Page,
} from '@influxdata/clockface';
import Otcls from './otcls';
import OtclEdit from './component/OtclEdit';
import OtclCreate from './component/OtclCreate';

import { OtclProvider } from './state';

const pageContentsClassName = `alerting-index alerting-index__${'check'}`;
const title = 'OpenTelemetry Collector';
const otclsPrefix = `/orgs/:orgID/otcls`;

const Header: React.FC = () => {
  return (
    <Page.Header fullWidth={true}>
      <Page.Title title={title} />
    </Page.Header>
  );
};

type OtclPageProps = {
  onCreate: () => void;
};

class OtclPage extends React.Component<OtclPageProps> {
  // todo: re-render after update

  public render() {
    const { onCreate } = this.props;

    return (
      <Page titleTag={title}>
        <Header />
        <Page.ControlBar fullWidth={true}>
          <Page.ControlBarRight>
            <Button
              size={ComponentSize.Small}
              icon={IconFont.Plus}
              color={ComponentColor.Primary}
              text={'Create Configuration'}
              onClick={onCreate}
            />
          </Page.ControlBarRight>
        </Page.ControlBar>

        <Page.Contents
          fullWidth={true}
          scrollable={false}
          className={pageContentsClassName}
        >
          <Otcls />
        </Page.Contents>
      </Page>
    );
  }
}

const Otcl: React.FC = () => {
  const { orgID } = useParams();
  const history = useHistory();
  const onCreate = useCallback(() => {
    history.push(`/orgs/${orgID}/otcls/new`);
  }, [orgID]);

  return (
    <OtclProvider orgID={orgID}>
      <OtclPage onCreate={onCreate} />
      <Switch>
        <Route path={`${otclsPrefix}/new`} component={OtclCreate} />
        <Route path={`${otclsPrefix}/:otclID`} component={OtclEdit} />
      </Switch>
    </OtclProvider>
  );
};

export default Otcl;
