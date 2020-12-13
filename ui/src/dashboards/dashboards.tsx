import React from "react";
import { Button, ComponentColor, IconFont, Page } from "@influxdata/clockface";

import SearchWidget from "../components/SearchWidget";
import DashboardCards from "./components/DashboardCards";
import { Dashboard } from "types";

const dashboards: Dashboard[] = [
  {
    id: "a",
    name: "name",
    desc: "desc",
    created: "111111",
    updated: "222222",
    orgID: "111111",
    panels: []
  }
];

const DashboardsIndex: React.FC = () => {
  return (
    <Page>
      <Page.Header fullWidth={false}>
        <Page.Title title={"Dashboards"} />
        {/* rateLimitAlert? */}
      </Page.Header>

      <Page.ControlBar fullWidth={false}>
        <Page.ControlBarLeft>
          <SearchWidget
            placeholder={"Filter dashboards..."}
            onSearch={v => console.log("v", v)}
          />


        </Page.ControlBarLeft>
        <Page.ControlBarRight>
          <Button
            text={"Add"}
            icon={IconFont.Plus}
            color={ComponentColor.Primary}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>

      <Page.Contents
        className="dashboards-index__page-contents"
        fullWidth={false}
        scrollable={true}
      >
        <DashboardCards dashboards={dashboards} />
      </Page.Contents>
    </Page>
  );
};

export default DashboardsIndex;
