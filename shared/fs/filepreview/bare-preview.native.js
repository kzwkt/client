// @flow
import * as React from 'react'
import * as Types from '../../constants/types/fs'
import * as Constants from '../../constants/fs'
import * as DispatchMappers from '../utils/dispatch-mappers'
import {globalStyles, globalColors, globalMargins, platformStyles} from '../../styles'
import {Box, ClickableBox, Text, Icon} from '../../common-adapters'
import {navigateUp} from '../../actions/route-tree'
import {connect, type Dispatch, type TypedState} from '../../util/container'
import {type BarePreviewProps} from './bare-preview'
import View from './view-container'

const mapStateToProps = (state: TypedState, {routeProps}: BarePreviewProps) => {
  const path = Types.stringToPath(routeProps.get('path', Constants.defaultPath))
  return {
    path,
    _pathItem: state.fs.pathItems.get(path) || Constants.makeUnknownPathItem(),
  }
}

const mapDispatchToProps = (dispatch: Dispatch) => ({
  onBack: () => dispatch(navigateUp()),
  _onAction: DispatchMappers.mapDispatchToOnAction(dispatch),
})

const mergeProps = ({path, _pathItem}, {onBack, _onAction}) => ({
  path,
  onBack,
  onAction: (event: SyntheticEvent<>) => _onAction(path, _pathItem.type, event),
})

type ConnectedBarePreviewProps = {
  path: Types.Path,

  onBack: () => void,
  onAction: (evt?: SyntheticEvent<>) => void,
}

const BarePreview = (props: ConnectedBarePreviewProps) => (
  <Box style={stylesContainer}>
    <Box style={stylesHeader}>
      <ClickableBox onClick={props.onBack}>
        <Text type="Body" style={stylesText}>
          Close
        </Text>
      </ClickableBox>
    </Box>
    <Box style={stylesContentContainer}>
      <View path={props.path} />
    </Box>
    <Box style={stylesFooter}>
      <Icon type="iconfont-ellipsis" onClick={props.onAction} color={globalColors.white} />
    </Box>
  </Box>
)

const stylesContainer = platformStyles({
  common: {
    ...globalStyles.flexBoxColumn,
    ...globalStyles.flexGrow,
    backgroundColor: globalColors.black,
  },
  isIOS: {
    marginTop: -20, // top status bar
  },
})

const stylesText = {
  color: globalColors.white,
}

const stylesHeader = {
  ...globalStyles.flexBoxRow,
  alignItems: 'center',
  height: 32,
  paddingLeft: globalMargins.tiny,
}

const stylesContentContainer = {
  ...globalStyles.flexGrow,
}

const stylesFooter = {
  ...globalStyles.flexBoxRow,
  alignItems: 'center',
  height: 32,
  paddingLeft: globalMargins.tiny,
}

export default connect(mapStateToProps, mapDispatchToProps, mergeProps)(BarePreview)
