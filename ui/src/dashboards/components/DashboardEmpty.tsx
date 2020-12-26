import React from 'react';
import { Button, ComponentColor, ComponentSize, EmptyState, IconFont } from '@influxdata/clockface';

const DashboardEmpty = () => {
  return (
    <div className={'dashboard-empty'}>
      <EmptyState size={ComponentSize.Large}>
        <div className={'dashboard-empty--graphic'}>
          <div className={'dashboard-empty--graphic-content'} />
        </div>

        <EmptyState.Text>
          This Dashboard doesn't have any <b>Cells</b>, let's create one!
        </EmptyState.Text>

        <Button
          text={'Add Cell'}
          size={ComponentSize.Medium}
          icon={IconFont.AddCell}
          color={ComponentColor.Primary}
          onClick={() => console.log('addCell')}
        />
      </EmptyState>
    </div>
  );
};

export default DashboardEmpty;
