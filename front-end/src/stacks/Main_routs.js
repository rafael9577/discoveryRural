import {createBottomTabNavigator} from "@react-navigation/bottom-tabs"

import home from "../screens/home"

const Tab = createBottomTabNavigator();

export default () => (
    
        <Tab.Navigator initialRouteName="home" screenOptions={{headerShown:false}}>
            <Tab.Screen name="home" component={home}  />
        </Tab.Navigator>
)
