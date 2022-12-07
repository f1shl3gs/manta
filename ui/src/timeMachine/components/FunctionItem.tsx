// Libraries
import React, {useRef} from 'react'

// Components
import {
  Appearance,
  Button,
  ComponentColor,
  ComponentSize,
  DapperScrollbars,
  Popover,
  PopoverInteraction,
  PopoverPosition,
} from '@influxdata/clockface'

// Constants
import {PromqlFunction} from 'src/shared/constants/promqlFunctions'

interface Props {
  fn: PromqlFunction
  onClickFn: (fn: PromqlFunction) => void
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
          <div
            className={'flux-function-docs'}
            data-testid={`flux-docs--${fn.name}`}
          >
            <DapperScrollbars autoHide={false}>
              <div className={'flux-toolbar--popover'}>
                {/* desc */}
                <article className={'flux-functions-toolbar--description'}>
                  <div className={'flux-function-docs--heading'}>
                    Description
                  </div>
                  <span style={{whiteSpace: 'pre-line'}}>{fn.desc}</span>
                </article>

                {/* arguments */}
                <article>
                  <div className={'flux-function-docs--heading'}>Arguments</div>
                  <div className={'flux-function-docs--snippet'}>
                    {fn.args.length === 0 ? (
                      <div className={'flux-function-docs--arguments'}>
                        None
                      </div>
                    ) : (
                      fn.args.map(arg => (
                        <div
                          key={arg.name}
                          className={'flux-function-docs--arguments'}
                        >
                          <span>{arg.name}:</span>
                          <span>{arg.type}</span>
                          <div>{arg.desc}</div>
                        </div>
                      ))
                    )}
                  </div>
                </article>

                {/* Example */}
                <article>
                  <div className={'flux-function-docs--heading'}>Example</div>
                  <div className={'flux-function-docs--snippet'}>
                    {fn.example}
                  </div>
                </article>

                {/* Links */}
                <p className={'tooltip--link'}>
                  Still have questions? Check out the{' '}
                  <a
                    target={'_blank'}
                    href={`https://prometheus.io/docs/prometheus/latest/querying/basics`}
                  >
                    PromQL Docs
                  </a>
                </p>
              </div>
            </DapperScrollbars>
          </div>
        )}
      />
      <dd
        // @ts-ignore
        ref={fnRef}
        data-testid={`flux--${testID}`}
        className={'flux-toolbar--list-item flux-toolbar--function'}
      >
        <code>{fn.name}</code>
        <Button
          testID={`flux--${testID}--inject`}
          text={'Inject'}
          onClick={_ => onClickFn(fn)}
          size={ComponentSize.ExtraSmall}
          className={'flux-toolbar--injector'}
          color={ComponentColor.Primary}
        />
      </dd>
    </>
  )
}

export default FunctionItem
