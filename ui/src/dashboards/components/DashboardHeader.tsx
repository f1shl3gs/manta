import React, { useCallback } from 'react';
import { Button, ComponentColor, IconFont, Page } from '@influxdata/clockface';
import RenamablePageTitle from 'components/RenamablePageTitle';
import { useFetch } from 'use-http';
import { useParams } from 'react-router-dom';
import PresentationModeToggle from './PresentationModeToggle';
import AutoRefreshDropdown from 'components/AutoRefreshDropdown/AutoRefreshDropdown';
import TimeRangeDropdown from './TimeRangeDropdown';
import { AutoRefreshDropdownOptions } from '../../constants/autoRefresh';

const DashboardHeader = () => {

  const { dashboardID } = useParams<{ dashboardID: string }>();
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
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          placeholder={'Name this dashboard'}
          name={'name this'}
          maxLength={90}
          onRename={(name) => console.log('name')}
        />
      </Page.Header>
      <Page.ControlBar fullWidth={true}>
        <Page.ControlBarLeft>
          <Button
            text={'Add Cell'}
            color={ComponentColor.Primary}
            icon={IconFont.AddCell}
            onClick={addCell}
          />
          <Button
            icon={IconFont.TextBlock}
            text="Add Note"
            onClick={() => console.log('add note')}
            testID="add-note--button"
          />
          <PresentationModeToggle />
        </Page.ControlBarLeft>

        <Page.ControlBarRight>
          <TimeRangeDropdown />
          <AutoRefreshDropdown options={AutoRefreshDropdownOptions} />
        </Page.ControlBarRight>
      </Page.ControlBar>
    </>
  );
};

export default DashboardHeader;
