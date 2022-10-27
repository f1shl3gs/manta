// Libraries
import React, {FC} from 'react'

interface Props {
  children: any
}

export const WizardFullScreen: FC<Props> = props => {
  const {children} = props

  return (
    <>
      <div className={'wizard--full-screen'}>
        {children}

        <div className={'wizard--credits'}>Powered by Manta</div>
      </div>
    </>
  )
}
