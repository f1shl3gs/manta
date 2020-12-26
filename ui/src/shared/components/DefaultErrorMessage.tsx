import React, {ComponentType} from 'react'

type ErrorMessageComponent = ComponentType<{error: Error}>

const DefaultErrorMessage: ErrorMessageComponent = () => {
  return (
    <p
      className="default-error-message"
      style={{display: 'flex', placeContent: 'center'}}
    >
      An InfluxDB error has occurred. Please report the issue&nbsp;
      <a href="https://github.com/influxdata/influxdb/issues">here</a>.
    </p>
  )
}

export default DefaultErrorMessage
