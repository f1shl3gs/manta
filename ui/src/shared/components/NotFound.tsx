import React from 'react'
import {AppWrapper, FunnelPage} from '@influxdata/clockface'

const NotFound: React.FC = () => (
  <AppWrapper type={'funnel'} className={'page-not-found'}>
    <FunnelPage enableGraphic={true} className={'page-not-found-funnel'}>
      <div>
        <h2 className={'page-not-found-content-highlight'}>
          404: Page Not Found
        </h2>

        <h4>
          We couldn't find the page you were looking for,
          <br />
          please refresh the page or check the URL and try again.
        </h4>
      </div>
    </FunnelPage>
  </AppWrapper>
)

export default NotFound
