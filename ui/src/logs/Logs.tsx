import React from "react";
import { Button, ComponentColor, Input, Page } from "@influxdata/clockface";
import { QueryProvider, useQuery } from "./useQuery";
import LogList from "./components/LogList";

const Title = "Logs";

const Header: React.FC = () => {
  return (
    <Page.Header fullWidth>
      <Page.Title title={Title} />
    </Page.Header>
  );
};

const Logs: React.FC = () => {
  const { query, setQuery } = useQuery();

  return (
    <Page titleTag={"logs"}>
      <Header />
      <Page.ControlBar fullWidth>
        <Page.ControlBarLeft>
          <Button text={"Signature"} />
        </Page.ControlBarLeft>
        <Page.ControlBarRight>
          <Input placeholder={"LogQL"} value={query} onChange={ev => setQuery(ev.target.value)} />
          <Button
            color={ComponentColor.Primary}
            text={"Query"} onClick={() => {
            console.log("query");
          }} />
        </Page.ControlBarRight>
      </Page.ControlBar>

      <Page.Contents
        fullWidth
        scrollable={true}
      >
        <LogList />
      </Page.Contents>
    </Page>
  );
};

const wrapped = () => (
  <QueryProvider>
    <Logs />
  </QueryProvider>
);

export default wrapped;
