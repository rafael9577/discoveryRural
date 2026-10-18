import { createStackNavigator } from '@react-navigation/stack'

import login from '../screens/login'
import singUp from '../screens/signup'
import tabsRouts from './TabsNavigator'

const stack = createStackNavigator();

export default () => (
    <stack.Navigator initialRouteName='login' screenOptions={{ headerShown: false }}>
        <stack.Screen name='login' component={login} />
        <stack.Screen name='singup' component={singUp} />
        <stack.Screen name='userRouts' component={tabsRouts}/>
    </stack.Navigator>
)