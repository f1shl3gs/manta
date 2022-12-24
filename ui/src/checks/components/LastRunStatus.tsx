// Libraries
import React, {FunctionComponent, useRef, useState} from 'react'
import classnames from 'classnames'

// Components
import {
  Appearance,
  ComponentColor,
  Icon,
  IconFont,
  Popover,
  PopoverInteraction,
} from '@influxdata/clockface'

interface Props {
  lastRunError?: string
  lastRunStatus: string
}

const LastRunTaskStatus: FunctionComponent<Props> = ({
  lastRunError,
  lastRunStatus,
}) => {
  const triggerRef = useRef<HTMLDivElement>(null)
  const [highlight, setHighlight] = useState(false)

  let color = ComponentColor.Success
  let icon = IconFont.CheckMark_New
  let text = 'Task ran successfully!'

  if (lastRunStatus === 'failed' || lastRunError !== undefined) {
    color = ComponentColor.Danger
    icon = IconFont.AlertTriangle
    text = lastRunError
  }

  if (lastRunStatus === 'cancel') {
    color = ComponentColor.Warning
    icon = IconFont.Remove_New
    text = 'Task Cancelled'
  }

  const statusClassName = classnames('last-run-task-status', {
    [`last-run-task-status__${color}`]: color,
    'last-run-task-status__highlight': highlight,
  })

  const popoverContents = () => (
    <>
      <h6>Last Run Status:</h6>
      <p>{text}</p>
    </>
  )

  return (
    <>
      <div
        data-testid="last-run-status--icon"
        ref={triggerRef}
        className={statusClassName}
      >
        <Icon glyph={icon} />
      </div>

      <Popover
        className={'last-run-task-status--popover'}
        enableDefaultStyles={false}
        color={color}
        appearance={Appearance.Outline}
        triggerRef={triggerRef}
        contents={popoverContents}
        showEvent={PopoverInteraction.Hover}
        hideEvent={PopoverInteraction.Hover}
        onShow={() => setHighlight(true)}
        onHide={() => setHighlight(false)}
      />
    </>
  )
}

export default LastRunTaskStatus
