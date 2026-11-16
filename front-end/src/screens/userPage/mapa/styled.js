import { StyleSheet, Dimensions } from 'react-native';

const { width } = Dimensions.get("window");

export const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0A5133', 
  },
  
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: 'transparent',
    zIndex: 2, 
  },
  headerSpacer: {
    width: 40, 
  },
  headerTitleBox: {
    flex: 1,
    alignItems: 'center',
  },
  headerTitleText: {
    fontSize: 22,
    fontWeight: 'bold',
    color: '#FFFFFF', 
  },
  menuButton: {
    width: 40,
    height: 40,
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 1,
    borderColor: '#1B5E20',
    borderRadius: 12,
    backgroundColor: '#2E7D32', 
  },
  menuButtonAtivo: {
    backgroundColor: '#1B5E20', 
  },

  /* MAPA estilos */
  mapContainer: {
    flex: 1,
    width: "95%",
    alignSelf: "center",
    marginVertical: 10,
    borderRadius: 20, 
    overflow: 'hidden', 
    borderWidth: 2,
    borderColor: '#1B5E20',
  },
  map: {
    width: '100%',
    height: '100%',
  },

  overlayMenu: {
    position: 'absolute',
    top: 10,
    right: 10,
    width: '75%', 
    backgroundColor: '#2E7D32', 
    borderRadius: 20,
    borderWidth: 2,
    borderColor: '#1B5E20',
    overflow: 'hidden',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
    elevation: 5,
  },
  overlayHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: '#1B5E20', 
    paddingHorizontal: 16,
    paddingVertical: 12,
  },
  overlayHeaderText: {
    color: '#FFFFFF',
    fontWeight: 'bold',
    fontSize: 14,
  },
  overlayBody: {
    backgroundColor: '#2E7D32', 
  },
  filterTitleContainer: {
    backgroundColor: 'rgba(0,0,0,0.2)',
    paddingVertical: 8,
    paddingHorizontal: 16,
  },
  filterTitleText: {
    fontWeight: 'bold',
    fontSize: 12,
    color: '#E0E0E0',
  },
  checkboxItemContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 12,
    paddingHorizontal: 16,
  },
  checkboxItemText: {
    fontSize: 14,
    color: '#FFFFFF',
    fontWeight: '500',
  },
  checkboxBox: {
    width: 24,
    height: 24,
    borderWidth: 2,
    borderColor: '#FFFFFF',
    borderRadius: 6,
    backgroundColor: 'transparent',
    justifyContent: 'center',
    alignItems: 'center',
  },
  separator: {
    height: 1,
    backgroundColor: '#1B5E20',
    marginHorizontal: 16,
  },
  
  /* EStilo do botão */
  avistamento: {
    backgroundColor: "#268596", 
    width: "95%",
    alignSelf: 'center',
    alignItems: 'center',
    paddingVertical: 15,
    borderRadius: 20,
    marginBottom: 20,
  },
  avistamentoText: {
    color: "#FFFFFF",
    fontWeight: "bold",
    fontSize: 16,
  },
  overlayHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: '#1B5E20', 
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderBottomWidth: 1, // Adicionado para separar
    borderBottomColor: '#2E7D32', // Adicionado para separar
  },
});