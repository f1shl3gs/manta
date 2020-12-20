import React, { useCallback } from "react";
import { Button, ComponentColor, IconFont, Page, SpinnerContainer, TechnoSpinner } from "@influxdata/clockface";
import { useFetch } from "use-http";
import { Route, Switch, useParams } from "react-router-dom";
import { useOrgID } from "../../shared/state/organization/organization";
import { Dashboard } from "../../types";
import RenamablePageTitle from "../../components/RenamablePageTitle";
import remoteDataState from "../../utils/rds";
import AutoRefreshDropdown from "../../components/AutoRefreshDropdown";
import { AutoRefreshOption } from "../../types/AutoRefresh";
import { DashboardProvider } from "../state/dashboard";
import Cells from "./Cells";
import ViewEditorOverlay from "./ViewEditorOverlay";

const autoRefreshDropdownOptions: AutoRefreshOption[] = [
  {
    label: "pause",
    seconds: 0
  },
  {
    label: "Last 5m",
    seconds: 5 * 60
  },
  {
    label: "Last 15m",
    seconds: 15 * 60
  },
  {
    label: "Last 30m",
    seconds: 30 * 60
  },
  {
    label: "Last 1h",
    seconds: 60 * 60
  },
  {
    label: "Last 3h",
    seconds: 3 * 60 * 60
  },
  {
    label: "Last 6h",
    seconds: 6 * 60 * 60
  }
];

const dashRoute = `/orgs/:orgID/dashboards/:dashboardID`

const DashboardPage: React.FC = () => {
  const { dashboardID } = useParams<{ dashboardID: string }>();
  const orgID = useOrgID();
  const { data, error, loading } = useFetch<Dashboard>(`/api/v1/dashboards/${dashboardID}?orgID=${orgID}`, {}, []);
  const rds = remoteDataState(loading, error);
  const { post } = useFetch(`/api/v1/dashboards/${dashboardID}/cells`, {});

  const addCell = useCallback(() => {
    return post({
      w: 4,
      h: 4,
      x: 0,
      y: 0
    });
  }, []);

  return (
    <>
      <Page>
        <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>

          <Page.Header fullWidth={true}>
            <RenamablePageTitle
              placeholder={"Name this dashboard"}
              name={data?.name || ""}
              maxLength={90}
              onRename={name => console.log("name")}
            />
          </Page.Header>

          <Page.ControlBar fullWidth={true}>
            <Page.ControlBarLeft>
              <Button
                text={"Add Cell"}
                color={ComponentColor.Primary}
                icon={IconFont.AddCell}
                onClick={addCell}
              />
            </Page.ControlBarLeft>

            <Page.ControlBarRight>
              <AutoRefreshDropdown options={autoRefreshDropdownOptions} />
            </Page.ControlBarRight>

          </Page.ControlBar>

          <Page.Contents>
            <Cells />
          </Page.Contents>

        </SpinnerContainer>
      </Page>

      <Switch>
        <Route path={`${dashRoute}/cells/:cellID/edit`} component={ViewEditorOverlay} />
      </Switch>
    </>
  );
};

const wrapper = () => (
  <DashboardProvider>
    <DashboardPage />
  </DashboardProvider>
);

export default wrapper;