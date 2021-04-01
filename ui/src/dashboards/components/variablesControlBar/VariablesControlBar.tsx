// Libraries
import React, {useCallback, useEffect, useState} from 'react'
import classnames from 'classnames'

// Components
import {
  ComponentSize,
  EmptyState,
  RemoteDataState,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'

// Hooks
import {usePresentationMode} from '../../../shared/usePresentationMode'
import ErrorBoundary from '../../../shared/components/ErrorBoundary'
import DraggableDropdown from './DraggableDropdown'

import withDragDropContext from '../../../shared/withDragDropContext'

const VariablesControlBar: React.FC = () => {
  const {inPresentationMode} = usePresentationMode()
  const [variables, setVariables] = useState([
    {
      id: 'foo',
      name: 'Foo',
    },
    {
      id: 'bar',
      name: 'Bar',
    },
  ])

  const emptyBar = () => {
    return (
      <EmptyState
        size={ComponentSize.ExtraSmall}
        className={'variables-control-bar--empty'}
      >
        <EmptyState.Text>
          This dashboard doesn't have any cells with defined variables
        </EmptyState.Text>
      </EmptyState>
    )
  }

  const handleMoveDropdown = useCallback(
    (oldIndex: number, newIndex: number) => {
      const next = [...variables]
      const tmp = next[oldIndex]
      next[oldIndex] = next[newIndex]
      next[newIndex] = tmp

      setVariables(next)
    },
    [variables, setVariables]
  )

  return (
    <div
      className={classnames('variables-control-bar', {
        'presentation-mode': inPresentationMode,
      })}
    >
      <SpinnerContainer
        className={'variables-spinner-container'}
        loading={RemoteDataState.Done}
        spinnerComponent={<TechnoSpinner diameterPixels={50} />}
      >
        {variables.length === 0 ? (
          emptyBar()
        ) : (
          <div className={'variables-control-bar--full'}>
            {variables.map((v, i) => (
              <ErrorBoundary key={v.id}>
                <DraggableDropdown
                  name={v.name}
                  id={v.id}
                  index={i}
                  moveDropdown={handleMoveDropdown}
                />
              </ErrorBoundary>
            ))}
          </div>
        )}
      </SpinnerContainer>
    </div>
  )
}

export default withDragDropContext(VariablesControlBar)
