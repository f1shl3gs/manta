// Libraries
import React, {FunctionComponent} from 'react'
import classnames from 'classnames'

interface Props {
  className?: string
  testID?: string
  widthPixels?: number
  children: JSX.Element | JSX.Element[]
}

const BuilderCard: FunctionComponent<Props> = ({
  testID,
  children,
  className,
  widthPixels = 228,
}) => {
  const style = {flex: `0 0 ${widthPixels}px`}
  const classname = classnames('builder-card', {[`${className}`]: className})

  return (
    <div className={classname} style={style} data-testid={testID}>
      {children}
    </div>
  )
}

export default BuilderCard
