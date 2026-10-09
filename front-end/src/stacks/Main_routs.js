import { createStackNavigator } from '@react-navigation/stack'

import login from '../screens/login'
import tabsRouts from './tabs_routs'

const stack = createStackNavigator();

export default () => (
    <stack.Navigator initialRouteName='login' screenOptions={{ headerShown: false }}>
        <stack.Screen name='login' component={login} />
        <stack.Screen name='userRouts' component={tabsRouts}/>
    </stack.Navigator>
)