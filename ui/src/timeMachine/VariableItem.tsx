import React, {useRef} from 'react'
import {Variable} from 'src/types/Variable'
import {
  Appearance,
  Button,
  ComponentColor,
  ComponentSize,
  Panel,
  Popover,
  PopoverInteraction,
  PopoverPosition,
} from '@influxdata/clockface'

interface Props {
  variable: Variable
  testID?: string
}

const VariableItem: React.FC<Props> = props => {
  const {variable, testID = ''} = props
  const trigger = useRef<HTMLDivElement>(null)

  return (
    <>
      <div
        className={'flux-toolbar--list-item flux-toolbar--variable'}
        ref={trigger}
        data-testid={`variable--${testID}`}
      >
        <code data-testid={`variable-name--${testID}`}>{variable.name}</code>
        <Button
          testID={`variable--${testID}--inject`}
          text={'Inject'}
          onClick={() => console.log('inject')}
          size={ComponentSize.ExtraSmall}
          className={'flux-toolbar--injector'}
          color={ComponentColor.Success}
        />
      </div>

      <Popover
        appearance={Appearance.Outline}
        position={PopoverPosition.ToTheLeft}
        triggerRef={trigger}
        showEvent={PopoverInteraction.Hover}
        hideEvent={PopoverInteraction.Hover}
        color={ComponentColor.Success}
        distanceFromTrigger={8}
        testID={'toolbar-popover'}
        enableDefaultStyles={false}
        contents={() => (
          <Panel>
            <Panel.Header>
              <h5>{variable.name}</h5>
            </Panel.Header>

            <Panel.Body size={ComponentSize.Small}>{variable.desc}</Panel.Body>
          </Panel>
        )}
      />
    </>
  )
}

export default VariableItem
