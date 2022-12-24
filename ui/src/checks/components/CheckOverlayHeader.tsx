// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ComponentColor,
  ComponentSize,
  FlexBox,
  IconFont,
  Page,
  SquareButton,
} from '@influxdata/clockface'
import RenamablePageTitle from 'src/shared/components/RenamablePageTitle'

interface Props {
  name: string
  onSetName: (name: string) => void
  onCancel: () => void
  onSave: () => void
}

const CheckOverlayHeader: FunctionComponent<Props> = ({
  name,
  onSetName,
  onCancel,
  onSave,
}) => {
  return (
    <>
      <Page.Header fullWidth={true}>
        <RenamablePageTitle
          onRename={onSetName}
          name={name}
          placeholder={'Name this check'}
          maxLength={64}
        />

        <FlexBox margin={ComponentSize.Small}>
          <SquareButton
            icon={IconFont.Remove_New}
            onClick={onCancel}
            size={ComponentSize.Small}
            testID={'check-overlay-cancel'}
          />
          <SquareButton
            icon={IconFont.CheckMark_New}
            color={ComponentColor.Primary}
            onClick={onSave}
            size={ComponentSize.Small}
            testID={'check-overlay-save'}
          />
        </FlexBox>
      </Page.Header>
    </>
  )
}

export default CheckOverlayHeader
