// Libraries
import React, {FunctionComponent} from 'react'

// Components
import DashedButton from 'src/shared/components/dashed_button/DashedButton'
import ThresholdStatement from 'src/checks/components/checkBuilder/ThresholdStatement'
import ThresholdRangeInput from 'src/checks/components/checkBuilder/ThresholdRangeInput'
import ThresholdValueInput from 'src/checks/components/checkBuilder/ThresholdValueInput'

// Types
import {
  Condition,
  ConditionStatus,
  Threshold,
  ThresholdType,
} from 'src/types/checks'
import {ComponentSize} from '@influxdata/clockface'

// Constants
import {LEVEL_COMPONENT_COLORS} from 'src/checks/constants'

// Actions
import {setCondition, removeCondition} from 'src/checks/actions/builder'
import {connect, ConnectedProps} from 'react-redux'

const mdtp = {
  removeCondition,
  setCondition,
}

interface OwnProps {
  status: ConditionStatus
  threshold?: Threshold
}

const connector = connect(null, mdtp)
type Props = OwnProps & ConnectedProps<typeof connector>

const ConditionCard: FunctionComponent<Props> = ({
  status,
  threshold,
  setCondition,
  removeCondition,
}) => {
  const handleAdd = () => {
    const condition: Condition = {
      status,
      pending: '1m',
      threshold: {
        type: 'gt',
        value: 100,
      },
    }

    setCondition(condition)
  }

  const handleChangeType = (type: ThresholdType) => {
    if (type === 'inside' || type === 'outside') {
      setCondition({
        status,
        pending: '0s',
        threshold: {
          type,
          min: 0,
          max: 100,
        },
      })
    } else {
      setCondition({
        status,
        pending: '0s',
        threshold: {
          type,
          value: 0,
        },
      })
    }
  }

  const handleDelete = () => {
    console.log('delete')

    removeCondition(status)
  }

  if (!threshold) {
    return (
      <DashedButton
        text={`+ ${status}`}
        color={LEVEL_COMPONENT_COLORS[status]}
        size={ComponentSize.Large}
        testID={`add-threshold-condition-${status}`}
        onClick={handleAdd}
      />
    )
  }

  const inner =
    threshold.type === 'inside' || threshold.type === 'outside' ? (
      <ThresholdRangeInput
        threshold={threshold}
        onChange={(min, max) => console.log('onchange', min, max)}
      />
    ) : (
      <ThresholdValueInput
        threshold={threshold}
        onChange={value => console.log('value', value)}
      />
    )

  return (
    <ThresholdStatement
      status={status}
      threshold={threshold}
      onChangeType={handleChangeType}
      onDelete={handleDelete}
    >
      {inner}
    </ThresholdStatement>
  )
}

export default connector(ConditionCard)
