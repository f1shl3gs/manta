// Libraries
import React, {useRef} from 'react'

// Components
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

interface FunctionTerm {
  label: string
  detail: string
  info: string
  type: string
}

interface Props {
  fn: FunctionTerm
  onClickFn: (fn: FunctionTerm) => void
  testID?: string
}

const FunctionItem: React.FC<Props> = props => {
  const {fn, onClickFn, testID = 'flux-function'} = props
  const fnRef = useRef<HTMLDListElement>()

  return (
    <>
      <Popover
        appearance={Appearance.Outline}
        enableDefaultStyles={false}
        position={PopoverPosition.ToTheLeft}
        triggerRef={fnRef}
        showEvent={PopoverInteraction.Hover}
        hideEvent={PopoverInteraction.Hover}
        distanceFromTrigger={8}
        testID={'toolbar-popover'}
        contents={() => (
          <Panel>
            <Panel.Header>
              <h5>{fn.label}</h5>
            </Panel.Header>

            <Panel.Body size={ComponentSize.Small}>
              <p>{fn.type}</p>
              <p>{fn.info}</p>
              <p>{fn.detail}</p>
            </Panel.Body>
          </Panel>
        )}
      />
      <dd
        // @ts-ignore
        ref={fnRef}
        data-testid={`flux--${testID}`}
        className={'flux-toolbar--list-item flux-toolbar--function'}
      >
        <code>{fn.label}</code>
        <Button
          testID={`flux--${testID}--inject`}
          text={'Inject'}
          onClick={ev => onClickFn(fn)}
          size={ComponentSize.ExtraSmall}
          className={'flux-toolbar--injector'}
          color={ComponentColor.Primary}
        />
      </dd>
    </>
  )
}

export default FunctionItem
