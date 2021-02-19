// Libraries
import React, {useRef, useState} from 'react'
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
  lastRunStatus: string
  lastRunError?: string
}

const LastRunStatus: React.FC<Props> = props => {
  const {lastRunError, lastRunStatus} = props
  const triggerRef = useRef<HTMLDivElement>(null)
  const [highlight, setHighlight] = useState<boolean>(false)

  let color = ComponentColor.Success
  let icon = IconFont.Checkmark
  let text = 'Task ran successfully!'

  const statusClassName = classnames('last-run-task-status', {
    [`last-run-task-status__${color}`]: color,
    'last-run-task-status__highlight': highlight,
  })

  if (lastRunStatus === 'failed' || lastRunError !== undefined) {
    color = ComponentColor.Danger
    icon = IconFont.AlertTriangle
    text = lastRunError || 'Unknown'
  }

  if (lastRunStatus === 'cancel') {
    color = ComponentColor.Warning
    icon = IconFont.Remove
    text = 'Task Cancelled'
  }

  const popoverContents = () => (
    <>
      <h6>Last Run Status:</h6>
      <p>{text}</p>
    </>
  )

  return (
    <>
      <div ref={triggerRef} className={statusClassName}>
        <Icon glyph={icon} />
      </div>

      <Popover
        className="last-run-task-status--popover"
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

export default LastRunStatus
