import React from 'react';

import { ComponentSize, IconFont, Notification } from '@influxdata/clockface';

import { useNotification } from '../shared/useNotification';

const Notifications = () => {
  const { notifications } = useNotification();

  return (
    notifications.map(item => {
      const {
        id,
        message,
        link
      } = item;

      return (
        <Notification
          key={id}
          id={id}
          icon={IconFont.Remove}
          duration={5}
          size={ComponentSize.Small}
          /*
              gradient={gradient}
              onTimeout={this.props.dismissNotification}
              onDismiss={this.props.dismissNotification}
              testID={`notification-${style}`}
         */

        >
          <span className={'notification--message'}>
            {message}
          </span>
        </Notification>
      );
    })
  );
};

export default Notifications;