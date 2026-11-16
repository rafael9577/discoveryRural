import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { useNavigation } from "@react-navigation/native";
import { Feather } from '@expo/vector-icons';
import { useSafeAreaInsets } from 'react-native-safe-area-context'; // Importado para corrigir o bug dos botões

import sair from '../screens/userPage/sair';
import busca from "../screens/userPage/informacao";
import emergencia from '../screens/userPage/emergencia';
import mapa from '../screens/userPage/mapa';

const tabs = createBottomTabNavigator();

const TelaVazia = () => null;

export default () => {
    const navigation = useNavigation();
    
    const insets = useSafeAreaInsets(); 

    return (
        <tabs.Navigator 
            initialRouteName='mapa' 
            screenOptions={({ route }) => ({
                headerShown: false,
                tabBarStyle: {
                    backgroundColor: '#FFFFFF', 
                    borderTopWidth: 0,
                    height: 60 + insets.bottom, 
                    paddingBottom: insets.bottom > 0 ? insets.bottom : 10, 
                    paddingTop: 10,                 
                    elevation: 10, 
                    shadowColor: '#000', 
                    shadowOffset: { width: 0, height: -3 },
                    shadowOpacity: 0.1,
                    shadowRadius: 4,
                },             
                tabBarActiveTintColor: '#0A5133', 
                tabBarInactiveTintColor: '#A0A0A0', 
                tabBarLabelStyle: {
                    fontSize: 12,
                    fontWeight: 'bold',
                    marginTop: 2,
                },

                tabBarIcon: ({ color, size }) => {
                    let iconName;
                    if (route.name === 'profile') {
                        iconName = 'user';
                    } else if (route.name === 'settings') {
                        iconName = 'settings';
                    } else if (route.name === 'busca') {
                        iconName = 'search';
                    } else if (route.name === 'sair') {
                        iconName = 'log-out';
                    } else if (route.name === 'prevenção') {
                        iconName = 'alert-octagon';
                    } else if (route.name === 'mapa') {
                        iconName = 'map';
                    } 
                    return <Feather name={iconName} size={size + 2} color={color} />; 
                },
            })}
        >
            <tabs.Screen name='prevenção' component={emergencia} />
            <tabs.Screen name='mapa' component={mapa}/>
            <tabs.Screen name='busca' component={busca}/>
            
            <tabs.Screen 
                name='sair' 
                component={TelaVazia} 
                listeners={{
                    tabPress: async (e) => {
                        e.preventDefault(); 
                        try {
                            await sair();     
                        } catch (error) {
                            console.log("Erro ao sair:", error);
                        } finally {

                            navigation.replace("login");
                        }
                    },
                }}
            />
        </tabs.Navigator>
    );
};