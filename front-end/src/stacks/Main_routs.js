import { createStackNavigator } from '@react-navigation/stack'

import login from '../screens/login'


const stack = createStackNavigator();

export default () => (
    <stack.Navigator initialRouteName='login' screenOptions={{ headerShown: false }}>
        <stack.Screen name='login' component={login} />
    </stack.Navigator>
)